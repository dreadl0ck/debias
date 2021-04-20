package debias

import (
	"bytes"
	"math"
)

func setBit(n byte, pos uint) byte {
	n |= 1 << pos
	return n
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	var (
		padding = blockSize - len(ciphertext)%blockSize
		padText = bytes.Repeat([]byte{byte(padding)}, padding)
	)
	return append(ciphertext, padText...)
}

func ShannonEntropy(data []byte) int {

	var (
		freqs = make(map[byte]float64)
		sum   float64
	)

	// track byte frequencies
	for _, i := range data {
		freqs[i]++
	}

	for _, v := range freqs {
		freq := v / float64(len(data))
		sum += freq * math.Log2(freq)
	}

	return int(math.Ceil(sum*-1)) * len(data)
}
