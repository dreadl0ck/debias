package debias

import (
	"bytes"
)

func setBit(n byte, pos uint) byte {
	n |= 1 << pos
	return n
}


func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	var (
		padding = blockSize - len(ciphertext) % blockSize
		padText = bytes.Repeat([]byte{byte(padding)}, padding)
	)
	return append(ciphertext, padText...)
}