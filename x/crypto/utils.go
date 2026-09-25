// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package crypto

import (
	"bytes"
	"errors"
)

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, errors.New(invalidBlockSize)
	}
	if len(b) == 0 {
		return nil, errors.New(invalidPKCS7Data)
	}
	n := blockSize - (len(b) % blockSize)
	pb := make([]byte, len(b)+n)
	copy(pb, b)
	copy(pb[len(b):], bytes.Repeat([]byte{byte(n)}, n))
	return pb, nil
}

// pkcs7UnPad validates and unPads data from the given bytes slice.
// The returned value will be 1 to n bytes smaller depending on the
// amount of padding, where n is the block size.
func pkcs7UnPad(b []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, errors.New(invalidBlockSize)
	}
	if len(b) == 0 {
		return nil, errors.New(invalidPKCS7Data)
	}
	if len(b)%blockSize != 0 {
		return nil, errors.New(invalidPKCS7Padding)
	}
	c := b[len(b)-1]
	n := int(c)
	if n == 0 || n > len(b) {
		return nil, errors.New(invalidPKCS7Padding)
	}
	for i := 0; i < n; i++ {
		if b[len(b)-n+i] != c {
			return nil, errors.New(invalidPKCS7Padding)
		}
	}
	return b[:len(b)-n], nil
}
