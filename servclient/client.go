package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"log"
	"net"
)

func encryptBlock(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()
	data = pkcs7Padding(data, blockSize)
	encrypted := make([]byte, len(data))

	iv := make([]byte, blockSize)
	_, err = rand.Read(iv)
	if err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(encrypted, data)

	return append(iv, encrypted...), nil
}

func decryptBlock(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < block.BlockSize() {
		return nil, fmt.Errorf("недостаточно данных для дешифрования")
	}

	iv := data[:block.BlockSize()]
	ciphertext := data[block.BlockSize():]

	mode := cipher.NewCBCDecrypter(block, iv)

	plainText := make([]byte, len(ciphertext))
	mode.CryptBlocks(plainText, ciphertext)

	plainText = pkcs7Unpadding(plainText)

	return plainText, nil
}

func pkcs7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func pkcs7Unpadding(data []byte) []byte {
	padding := int(data[len(data)-1])
	return data[:len(data)-padding]
}

func main() {
	key := []byte("1234567890abcdef1234567890abcdef")
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		log.Fatal("Ошибка при подключении к серверу:", err)
	}
	defer conn.Close()

	message := "DSVPN test"
	encryptedMessage, err := encryptBlock([]byte(message), key)
	if err != nil {
		log.Fatal("Ошибка при шифровании сообщения:", err)
	}

	_, err = conn.Write(encryptedMessage)
	if err != nil {
		log.Fatal("Ошибка при отправке данных на сервер:", err)
	}
	fmt.Println("Отправлено зашифрованное сообщение на сервер.")

	encryptedResponse := make([]byte, 1024)
	n, err := conn.Read(encryptedResponse)
	if err != nil {
		log.Fatal("Ошибка при получении данных с сервера:", err)
	}

	response, err := decryptBlock(encryptedResponse[:n], key)
	if err != nil {
		log.Fatal("Ошибка при дешифровании ответа:", err)
	}

	fmt.Printf("Ответ от сервера: %s\n", response)
}
