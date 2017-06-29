package encrypt

import "testing"

func TestCBFEncrypt(t *testing.T) {
	sourceFile := "/Users/hyh/Downloads/encrypt/test/0AACC39503EE31FC5BAF6A4DB27F2EC2_110080_120080_199_v02_mp4.ts"
	encryptFile := "/Users/hyh/Downloads/encrypt/test/encrypt.ts"
	decryptFile := "/Users/hyh/Downloads/encrypt/test/decrypt.ts"
	CBCEncryptFile(sourceFile, encryptFile, "7UhREDHUzVcRprq0", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	CBCDecryptFile(encryptFile, decryptFile, "7UhREDHUzVcRprq0", []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
}
