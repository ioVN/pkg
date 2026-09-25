// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

// AESGenerate256Key : generate an 256 bits AES key.
func AESGenerate256Key() (key []byte, err error) {
	b := make([]byte, 32)
	_, err = rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return key, err
	}
	return b, nil
}

// AESEncryptCBC padding with pkcs7 before encrypt with CBC mode.
func AESEncryptCBC(keyByte []byte, plaintext []byte) ([]byte, error) {
	plaintext, err := pkcs7Pad(plaintext, aes.BlockSize)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyByte)
	if err != nil {
		return nil, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)

	return ciphertext, nil
}

// AESDecryptCBC decrypt aes in CBC mode and remove pkcs7 padding.
func AESDecryptCBC(keySt []byte, ciphertext []byte) (decrypted []byte, err error) {
	block, err := aes.NewCipher(keySt)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New(invalidCiphertext)
		//panic()
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = pkcs7UnPad(ciphertext, aes.BlockSize)
	if err != nil {
		return nil, err
	}
	return ciphertext, nil
}
