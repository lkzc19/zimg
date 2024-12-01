package utils

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

// Boom 艺术就是爆炸
func Boom(err error) {
	if err != nil {
		fmt.Println("error:")
		fmt.Println("\t", err)
		os.Exit(-1)
	}
}

func GetBytes(filepath string) []byte {
	file, err := os.Open(filepath)
	Boom(err)
	defer file.Close()

	data, err := io.ReadAll(file)
	Boom(err)

	return data
}

func ToBase64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

func ToMd5(data []byte) string {
	hash := md5.Sum(data)
	return hex.EncodeToString(hash[:])
}
