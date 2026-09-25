// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package crypto

import (
	"os"
	"testing"
)

func TestOpenSSL_GenRSA(t *testing.T) {
	ssl := &OpenSSL{}
	// Generate on RAM
	if err := ssl.GenRSA(RSA2048Bits); err != nil {
		t.Errorf("GenRSA() error = %v", err)
	}
	// Generate to save file
	args := []string{"-out output/private.pem"}
	if err := ssl.GenRSA(RSA2048Bits, args...); err != nil {
		t.Errorf("GenRSA() error = %v", err)
	}
}

func TestOpenSSL_ExportP12(t *testing.T) {
	ssl := &OpenSSL{}
	// Generate to save file
	args := []string{"-out output/private.pem"}
	if err := ssl.GenRSA(RSA2048Bits, args...); err != nil {
		t.Errorf("GenRSA() error = %v", err)
	}
	//
	if err := ssl.ExportP12("123456", "output/private.p12"); err != nil {
		t.Errorf("ExportP12() error = %v", err)
	}
}

func TestOpenSSL_ImportP12(t *testing.T) {
	ssl := &OpenSSL{}
	derBytes, err := os.ReadFile("output/private.p12")
	if err != nil {
		t.Errorf("ReadFile() error = %v", err)
	}
	if err := ssl.ImportP12(derBytes, "123456"); err != nil {
		t.Errorf("ImportP12() error = %v", err)
	}
}
