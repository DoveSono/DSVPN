// gost/gost.go
package gost

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
)

// Генерация 256-битного ключа для ГОСТ 28147-89
func GenerateGostKey() (string, error) {
	key := make([]byte, 32) // 256 бит
	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(key), nil
}

// Генерация S-блока ГОСТ 28147-89 в формате 8x16
func GenerateGostSBlock() string {
	sBlock := ""
	for i := 0; i < 8; i++ {
		for j := 0; j < 16; j++ {
			num, _ := rand.Int(rand.Reader, big.NewInt(16))
			sBlock += fmt.Sprintf("%d", num.Int64())
			if j < 15 {
				sBlock += ","
			}
		}
		sBlock += "\n"
	}
	return sBlock
}
