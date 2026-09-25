// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package crypto

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"
	"log"
	"os"

	"golang.org/x/crypto/pkcs12"
)

// Keypair :
type Keypair struct {
	PrivateKey []byte
	PublicKey  []byte
}

// RSAGenerateKey : generate publicKey (DER-encoded PKIX format), privateKey (PKCS#8 encoded form, see RFC 5208).
func _() (kp *Keypair, err error) {
	priKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	privateKey, err := x509.MarshalPKCS8PrivateKey(priKey)
	if err != nil {
		return nil, err
	}
	publicKey, err := x509.MarshalPKIXPublicKey(&priKey.PublicKey)
	if err != nil {
		return nil, err
	}
	kp = &Keypair{privateKey, publicKey}
	return kp, nil
}

// RSAEncryptOAEP : encrypts the given message with RSA-OAEP, publicKey in DER-encoded PKIX format.
func RSAEncryptOAEP(publicKey, data []byte) (output []byte, err error) {
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		log.Println("EncryptRSA", err)
		return output, err
	}
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return output, errors.New("pub-key isn't a rsa pubkey")
	}

	label := []byte("")
	output, err = rsa.EncryptOAEP(sha1.New(), rand.Reader, pubKey, data, label)
	if err != nil {
		return output, err
	}
	return output, nil
}

// RSADecryptOAEP decrypts ciphertext using RSA-OAEP, privateKey in PKCS#8 encoded form.
func RSADecryptOAEP(ciphertext, privateKey []byte) (output []byte, err error) {
	p, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	pk, ok := p.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("private key wrong type")
	}
	label := []byte("")
	output, err = rsa.DecryptOAEP(sha1.New(), rand.Reader, pk, ciphertext, label)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return output, nil
}

// RSADecodePKCS12PriKey load rsa private key from p12 file
func _(privateKeyPath string, password string) (*rsa.PrivateKey, error) {
	pfxData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}
	priKey, _, err := pkcs12.Decode(pfxData, password)
	if err != nil {
		return nil, err
	}
	rsaPriKey, ok := priKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New(invalidPrivateKey)
	}
	return rsaPriKey, nil
}

// RSASignMessage sign message, return base64 signature
func RSASignMessage(message []byte, signKey []byte) (signature string, err error) {
	priKey, err := x509.ParsePKCS8PrivateKey(signKey)
	if err != nil {
		return "", err
	}
	privily, ok := priKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("private-key wrong type")
	}
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	sig, err := rsa.SignPKCS1v15(rand.Reader, privily, crypto.SHA256, d)
	if err != nil {
		return "", err
	}
	signature = base64.StdEncoding.EncodeToString(sig)
	return signature, nil
}

// RSAVerifySignature verify signature, return nil is signature valid, sig is signature base64
// sig base64 encoding, message raw message
func RSAVerifySignature(sig string, message []byte, verifyKey []byte) error {
	pub, err := x509.ParsePKIXPublicKey(verifyKey)
	if err != nil {
		log.Println("EncryptRSA", err)
		return err
	}
	pubKey, ok := pub.(*rsa.PublicKey)
	if !ok {
		return errors.New("pubkey isn't a rsa pubkey")
	}
	h := sha256.New()
	h.Write(message)
	d := h.Sum(nil)
	sigBytes, err := base64.StdEncoding.DecodeString(sig)
	if err != nil {
		return err
	}
	return rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, d, sigBytes)
}
