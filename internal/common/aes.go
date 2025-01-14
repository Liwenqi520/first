package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

func Encrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	plainBytes := []byte(plaintext)
	// 对于CBC模式，需要使用PKCS#7填充plaintext到blocksize的整数倍
	plainBytes = Pad(plainBytes, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plainBytes))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], plainBytes)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(key []byte, ct string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(ct)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", err
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)
	// 删除PKCS#7填充
	plaintext := Unpad(ciphertext)
	return string(plaintext), nil
}

// pad 使用PKCS#7标准填充数据
func Pad(buf []byte, blockSize int) []byte {
	padding := blockSize - (len(buf) % blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(buf, padtext...)
}

// unpad 删除PKCS#7填充的字节
func Unpad(buf []byte) []byte {
	length := len(buf)
	if length == 0 {
		return buf
	}
	unpadding := int(buf[length-1])
	return buf[:length-unpadding]
}
