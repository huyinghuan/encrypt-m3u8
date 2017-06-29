package encrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encry/utils"
	"io"
	"io/ioutil"
	"os"
)

func CBCEncryptStream(content []byte, aeskey string, iv []byte) ([]byte, error) {
	key := []byte(aeskey)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	content, err = utils.PKCS7Padding(content, block.BlockSize())
	if err != nil {
		return nil, err
	}
	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return nil, err
	// }
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(content, content)
	return content, nil
}

//CBCEncryptFile  aes-cbc-128 加密文件
// source 源文件路径
// dist  加密后输出文件路径
func CBCEncryptFile(source string, dist string, aeskey string, iv []byte) error {
	plaintext, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}
	key := []byte(aeskey)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	plaintext, err = utils.PKCS7Padding(plaintext, block.BlockSize())
	if err != nil {
		return err
	}
	// if _, err := io.ReadFull(rand.Reader, iv); err != nil {
	// 	return err
	// }
	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(plaintext, plaintext)
	f, err := os.Create(dist)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(plaintext))
	if err != nil {
		return err
	}
	return nil
}

//CBCDecryptFile  aes-cbc-128 解密文件
// source 源文件路径
// dist  加密后输出文件路径
func CBCDecryptFile(source string, dist string, aesKey string, iv []byte) error {
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
	bm := cipher.NewCBCDecrypter(block, iv)
	bm.CryptBlocks(ciphertext, ciphertext)
	ciphertext, err = utils.PKCS7Trimming(ciphertext, aes.BlockSize)
	if err != nil {
		return err
	}
	f, err := os.Create(dist)
	if err != nil {
		return err
	}
	_, err = io.Copy(f, bytes.NewReader(ciphertext))
	if err != nil {
		return err
	}
	return nil
}
