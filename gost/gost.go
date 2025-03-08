package gost

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
)

func addPadding(data []byte) []byte {
	padding := aes.BlockSize - len(data)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func removePadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

func GenerateGostKey() (string, error) {
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

func GostEncryptBlock(block, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("ключ должен быть 256 бит")
	}

	block = addPadding(block)

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := make([]byte, aes.BlockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(blockCipher, iv)

	encryptedBlock := make([]byte, len(block))
	mode.CryptBlocks(encryptedBlock, block)

	return append(iv, encryptedBlock...), nil
}

func GostDecryptBlock(block, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("ключ должен быть 256 бит")
	}

	if len(block) < aes.BlockSize {
		return nil, errors.New("недостаточно данных для дешифрования")
	}

	iv := block[:aes.BlockSize]
	encryptedBlock := block[aes.BlockSize:]

	blockCipher, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCDecrypter(blockCipher, iv)

	decryptedBlock := make([]byte, len(encryptedBlock))
	mode.CryptBlocks(decryptedBlock, encryptedBlock)

	decryptedBlock = removePadding(decryptedBlock)

	return decryptedBlock, nil
}
