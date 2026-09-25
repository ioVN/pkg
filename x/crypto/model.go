// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package crypto

const (
	// invalidBlockSize indicates hash blockSize <= 0.
	invalidBlockSize string = "invalid blockSize"
	// invalidPKCS7Data indicates bad input to PKCS7 pad or unPad.
	invalidPKCS7Data string = "invalid PKCS7 data (empty or not padded)"
	// invalidPKCS7Padding indicates PKCS7 unPad fails to bad input.
	invalidPKCS7Padding = "invalid padding on input"

	// invalidCiphertext must bigger blockSize
	invalidCiphertext string = "ciphertext too short"

	// invalidPrivateKey p12 file invalid
	invalidPrivateKey string = "invalid file"
)
