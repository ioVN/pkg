// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/ioVN/pkg/x/crypto/pkcs12"
)

const (
	RSA2048Bits uint32 = 1 << 11
)

var (
	CertOptDefault = CertOption{
		Issuer: pkix.Name{
			Country:            []string{"VN"},
			Organization:       []string{},
			OrganizationalUnit: []string{},
			CommonName:         "",
		},
		SubjectNameObject: pkix.Name{
			Country:            []string{"VN"},
			Organization:       []string{"ioVN Co., Ltd"},
			OrganizationalUnit: []string{},
			Locality:           []string{},
			Province:           []string{"Some-State"},
			CommonName:         "",
		},
		MakeExpires: func() time.Time {
			return time.Now().AddDate(1, 0, 0)
		},
	}
)

type CertOption struct {
	Issuer, SubjectNameObject pkix.Name
	MakeExpires               func() time.Time
	EmailAddresses            []string
}

type OpenSSL struct {
	private *rsa.PrivateKey
	cert    *x509.Certificate
}

/*
ImportP12 with private.p12

$ openssl pkcs12 -in private.p12 -clcerts -nokeys -out PublicKey.cer

$ openssl pkcs12 -in private.p12 -nodes -nocerts | openssl rsa -out PrivateKey.key
*/
func (ssl *OpenSSL) ImportP12(data []byte, pwd string) error {
	key, cert, err := pkcs12.Decode(data, pwd)
	if err != nil {
		return err
	}
	p, ok := key.(*rsa.PrivateKey)
	if ok {
		ssl.private = p
		ssl.cert = cert
	} else {
		return errors.New("private key missing")
	}
	return nil
}

/*
GenRSA : This command generates an RSA private key.

$ openssl genrsa -out private.pem 2048
*/
func (ssl *OpenSSL) GenRSA(bits uint32, args ...string) (err error) {
	ssl.private, err = rsa.GenerateKey(rand.Reader, int(bits))
	if err != nil {
		return err
	}
	for _, arg := range args {
		if strings.HasPrefix(arg, "-out ") {
			s := strings.Split(arg, "-out ")
			switch len(s) {
			case 0, 1:
				println("s[1]= NULL")
				continue
			default:
				println("s[1]='" + s[1] + "'")
				data := pem.EncodeToMemory(&pem.Block{
					Type:  "PRIVATE KEY",
					Bytes: x509.MarshalPKCS1PrivateKey(ssl.private),
				})
				if err := os.WriteFile(s[1], data, 0666); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

/*
MakeCertificate : Make *x509.Certificate value for OpenSSL object.

Note: Must call after GenRSA() function
*/
func (ssl *OpenSSL) MakeCertificate(opt *CertOption) error {
	if ssl.private == nil {
		return errors.New("must call after GenRSA() function")
	}
	var (
		certificateSerialNumber = time.Now().Unix()
	)
	template := x509.Certificate{
		SerialNumber:          big.NewInt(certificateSerialNumber),
		Issuer:                opt.Issuer,
		Subject:               opt.SubjectNameObject,
		NotBefore:             time.Now(),
		NotAfter:              opt.MakeExpires(), // Period too long
		SignatureAlgorithm:    x509.SHA256WithRSA,
		PublicKeyAlgorithm:    x509.RSA,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		EmailAddresses:        opt.EmailAddresses,
	}
	parseParent := func(p x509.Certificate) *x509.Certificate { p.Subject = template.Issuer; return &p }
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, parseParent(template),
		ssl.private.Public(), ssl.private)
	if err != nil {
		return err
	}
	ssl.cert, err = x509.ParseCertificate(derBytes)
	if err != nil {
		return err
	}
	return nil
}

/*
ExportP12 : save keyPair to pkcs12 type.

Note: Must call after GenRSA() or ImportP12() function
*/
func (ssl *OpenSSL) ExportP12(pwd string, output string) error {
	if ssl.private == nil {
		return errors.New("must call after GenRSA() function")
	}
	if ssl.cert == nil {
		if err := ssl.MakeCertificate(&CertOptDefault); err != nil {
			return err
		}
	}
	p12Bytes, err := pkcs12.Encode(rand.Reader, ssl.private, ssl.cert, nil, pwd)
	if err != nil {
		return err
	}
	if err := os.WriteFile(output, p12Bytes, 0644); err != nil {
		return err
	}
	return nil
}

/*
PrivateKey : Get *rsa.PrivateKey object.

Note: Must call after GenRSA() or ImportP12() function
*/
func (ssl *OpenSSL) PrivateKey() *rsa.PrivateKey {
	if ssl.private == nil {
		panic("must call after GenRSA() function")
	}
	return ssl.private
}

/*
Certificate : Get  *x509.Certificate object.

Note: Must call after  GenRSA() or ImportP12() function
*/
func (ssl *OpenSSL) Certificate() *x509.Certificate {
	if ssl.cert == nil {
		if err := ssl.MakeCertificate(&CertOptDefault); err != nil {
			panic("must call after GenRSA() function")
		}
	}
	return ssl.cert
}
