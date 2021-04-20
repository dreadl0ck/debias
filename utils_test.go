package debias_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
)

type bitString string

func (b bitString) bytes() []byte {
	var (
		out []byte
		str string
	)

	for i := len(b); i > 0; i -= 8 {
		if i-8 < 0 {
			str = string(b[0:i])
		} else {
			str = string(b[i-8 : i])
		}
		v, err := strconv.ParseUint(str, 2, 8)
		if err != nil {
			return nil
		}
		out = append([]byte{byte(v)}, out...)
	}

	return out
}

func (b bitString) hex() []string {
	var (
		out       []string
		byteSlice = b.bytes()
	)
	for _, b := range byteSlice {
		out = append(out, "0x"+hex.EncodeToString([]byte{b}))
	}
	return out
}

func TestBitStringConversion(t *testing.T) {
	if !bytes.Equal(bitString("00000000").bytes(), []byte{byte(0)}) {
		t.Fatal("incorrect conversion result: ", bitString("00000000").bytes(), " expected: ", []byte{byte(0)})
	}
	if !bytes.Equal(bitString("11111111").bytes(), []byte{byte(255)}) {
		t.Fatal("incorrect conversion result: ", bitString("11111111").bytes(), " expected: ", []byte{byte(255)})
	}
	if !bytes.Equal(bitString("10101010").bytes(), []byte{byte(170)}) {
		t.Fatal("incorrect conversion result: ", bitString("10101010").bytes(), " expected: ", []byte{byte(170)})
	}
	if !bytes.Equal(bitString("01010101").bytes(), []byte{byte(85)}) {
		t.Fatal("incorrect conversion result: ", bitString("01010101").bytes(), " expected: ", []byte{byte(85)})
	}
}

// Test for byte iteration implementation used during von neumann debiasing.
// We iterate over the bytes from left to right, that helped during manual creation of the unit tests.
// This test ensures the byte is iterated in the expected direction.
func TestIterateByte(t *testing.T) {

	var (
		b      = byte(216) // 11011000
		values = [][]int{
			{1, 1},
			{0, 1},
			{1, 0},
			{0, 0},
		}
		counter = 0
	)

	for j := 0; j < 8; j += 2 {
		ch := (b >> (7 - j)) & 0x01
		ch2 := (b >> (7 - (j + 1))) & 0x01

		if fmt.Sprint(values[counter][0]) != fmt.Sprint(ch) {
			t.Fatal("unexpected value")
		}
		if fmt.Sprint(values[counter][1]) != fmt.Sprint(ch2) {
			t.Fatal("unexpected value")
		}
		counter++
	}
}
