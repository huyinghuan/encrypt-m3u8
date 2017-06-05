package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encry/utils"
	"io"
	"io/ioutil"
	"os"
)

//CFBEncryptFile  aes-cbc-128 加密文件
// source 源文件路径
// dist  加密后输出文件路径
func CFBEncryptFile(source string, dist string, aeskey string) error {
	plaintext, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	// Byte array of the string
	key := []byte(aeskey)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	plaintext, err = utils.PKCS7Padding(plaintext, block.BlockSize())
	// The IV needs to be unique, but not secure. Therefore it's common to
	// include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	f, err := os.Create(dist)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	return err

}

func CFBDencryptFile(source string, dist string, aesKey string) error {
	ciphertext, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	// Key
	key := []byte(aesKey)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]
	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]
	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)
	// Decrypt bytes from ciphertext

	stream.XORKeyStream(ciphertext, ciphertext)
	// create a new file for saving the encrypted data.

	ciphertext, err = utils.PKCS7Trimming(ciphertext, aes.BlockSize)

	f, err := os.Create(dist)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	return err
}
