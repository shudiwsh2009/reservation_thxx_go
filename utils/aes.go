package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

/*
末尾有值为0的字节的明文 时候，该方法会有问题，确切的说zeropadding这种方法本身就有缺陷
*/
func AesDecrypt(cipherText, key, iv []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, fmt.Errorf("Empty input")
	}
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aesBlock.BlockSize() {
		return nil, fmt.Errorf("IV key should have the same size of BlockSize[%d], but it's length is %d", aesBlock.BlockSize(), len(iv))
	}
	if len(cipherText)%aesBlock.BlockSize() != 0 {
		return nil, fmt.Errorf("Cipher data should have size of multiple block size[%d]", aesBlock.BlockSize())
	}

	decryptor := cipher.NewCBCDecrypter(aesBlock, iv)
	decData := make([]byte, len(cipherText))
	decryptor.CryptBlocks(decData, cipherText)

	// we are using ZeroPadding from frontend
	lastZero := len(decData) - 1
	for i := len(decData) - 1; i >= 0; i -= 1 {
		if decData[i] != 0 {
			lastZero = i + 1
			break
		}
	}
	decData = decData[:lastZero]
	return decData, nil
}

func AesDecryptJson(cipherText, key, iv []byte, v interface{}) error {
	data, err := AesDecrypt(cipherText, key, iv)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

func AesEncrypt(input, key, iv []byte) ([]byte, error) {
	return aesEncryptWithPaddingFunc(input, key, iv, ZeroPadding)
}

type PaddingFunc func([]byte, int) []byte

func aesEncryptWithPaddingFunc(input, key, iv []byte, padding PaddingFunc) ([]byte, error) {
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aesBlock.BlockSize() {
		return nil, fmt.Errorf("IV key should have the same size of BlockSize[%d], but it's length is %d", aesBlock.BlockSize(), len(iv))
	}

	input = padding(input, aesBlock.BlockSize())
	encryptor := cipher.NewCBCEncrypter(aesBlock, iv)
	decData := make([]byte, len(input))
	encryptor.CryptBlocks(decData, input)
	return decData, nil
}

func AesEncryptWithPKCS5Padding(input, key, iv []byte) ([]byte, error) {
	return aesEncryptWithPaddingFunc(input, key, iv, PKCS5Padding)
}

func AesEncryptWithPKCS5PaddingAndUrlEncodeToString(input, key, iv []byte) (string, error) {
	encryptData, err := AesEncryptWithPKCS5Padding(input, key, iv)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encryptData), nil
}

func AesDecryptWithPKCS5Padding(cipherText, key, iv []byte) ([]byte, error) {
	if len(cipherText) == 0 {
		return nil, fmt.Errorf("Empty input")
	}
	aesBlock, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(iv) != aesBlock.BlockSize() {
		return nil, fmt.Errorf("IV key should have the same size of BlockSize[%d], but it's length is %d", aesBlock.BlockSize(), len(iv))
	}
	if len(cipherText)%aesBlock.BlockSize() != 0 {
		return nil, fmt.Errorf("Cipher data should have size of multiple block size[%d]", aesBlock.BlockSize())
	}

	decryptor := cipher.NewCBCDecrypter(aesBlock, iv)
	decData := make([]byte, len(cipherText))
	decryptor.CryptBlocks(decData, cipherText)

	return PKCS5UnPadding(decData)
}
