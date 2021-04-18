package debias_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dreadl0ck/debias"
)

type bitString string

func (b bitString) Bytes() []byte {
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

func (b bitString) Hex() []string {
    var out []string
    byteSlice := b.Bytes()
    for _, b := range byteSlice {
        out = append(out, "0x" + hex.EncodeToString([]byte{b}))
    }
    return out
}

func TestBitStringConversion(t *testing.T) {
	if !bytes.Equal(bitString("00000000").Bytes(), []byte{byte(0)}) {
		t.Fatal("incorrect conversion result: ", bitString("00000000").Bytes(), " expected: ", []byte{byte(0)})
	}
	if !bytes.Equal(bitString("11111111").Bytes(), []byte{byte(255)}) {
		t.Fatal("incorrect conversion result: ", bitString("11111111").Bytes(), " expected: ", []byte{byte(255)})
	}
	if !bytes.Equal(bitString("10101010").Bytes(), []byte{byte(170)}) {
		t.Fatal("incorrect conversion result: ", bitString("10101010").Bytes(), " expected: ", []byte{byte(170)})
	}
	if !bytes.Equal(bitString("01010101").Bytes(), []byte{byte(85)}) {
		t.Fatal("incorrect conversion result: ", bitString("01010101").Bytes(), " expected: ", []byte{byte(85)})
	}
}

// IMPORTANT: both input and output bit strings have to be a multiple of eight!
// An exception to this is an empty string for the output.
type test struct {
	in string
	out string
}

var tests = []test{
	{
		// 0 to 0000000000000000000 (steps of one) --> nothing
		in: "0000000000000000000",
		out: "",
	},
	{
		// 1 to 1111111111111111111 (steps of one) --> nothing
		in: "1111111111111111111",
		out: "",
	},
	{
		// 01 to 01010101 01010101 (steps of two) --> all zeros * 1/2 input length
		in: "0101010101010101",
		out: "00000000",
	},
	{
		// 10 to 10101010 10101010 (steps of two) --> all ones * 1/2 input length
		in: "1010101010101010",
		out: "11111111",
	},
	{
		in: "01101010110011111101100010100110",
		out: "0111011101000000", // contains padding with 6 trailing zeroes because Go's buffers can only hold full bytes.
	},
	{	
		in: "01001011001100100001011010100001",
		out: "0110011100000000", // seven trailing zeroes for padding
	},
	{	
		in: "00111110101000110000100001100011",
		out: "11110100", // two trailing zeroes for padding
	},
	
	{	
		in: "00111110101000110000100001100011",
		out: "11110100", // two trailing zeroes for padding
	},
	{	
		in: "01001000100100110001101001101101",
		out: "0110011010000000", // six trailing zeroes for padding
	},
	{	
		in: "00111111110011010011100011110011",
		out: "01000000", // six trailing zeroes for padding
	},
	{	
		in: "10110101011001111110110001010011",
		out: "1000101000000000", // six trailing zeroes for padding
	},
	{	
		in: "01100101100110010000101101010000",
		out: "0100101010000000", // five trailing zeroes for padding
	},
	{	
		in: "10011111010100011000010000110001",
		out: "10000100",
	},
	{	
		in: "01100100010010011000110100110110",
		out: "0100101001000000", // six trailing zeroes for padding
	},
	{	
		in: "10011111111001101001110001111001",
		out: "1010110010000000", // six trailing zeroes for padding
	},
}

func TestVonNeumann(t *testing.T) {
	for i, te := range tests {

		buf := debias.VonNeumann(bitString(te.in).Bytes())
				
		if !bytes.HasPrefix(buf.Bytes(), bitString(te.out).Bytes()) {

			fmt.Printf("buf.Bytes(): %08b\n", buf.Bytes())
			fmt.Printf("bitString(te.out).AsByteSlice(): %08b\n", bitString(te.out).Bytes())

			spew.Dump(tests[i])
			fmt.Println("input size", strconv.Itoa(len(te.in)) + " bits")
			fmt.Println("output size", strconv.Itoa(len(te.out)) + " bits")
			fmt.Println("======================= DEBUG ========================")
			debias.VonNeumann(bitString(te.in).Bytes())
			fmt.Println("======================= DEBUG ========================")
			t.Fatal("test #", i, ": expected " + te.out + " (" + strconv.Itoa(len(te.out)) + " bits) but got " + fmt.Sprintf("%08b", buf.Bytes()) + " (" + strconv.Itoa(len(buf.Bytes()) * 8) + " bits)")
		}
	}
}

func TestIterateByte(t *testing.T) {
	b := byte(216) // 11011000

	values := [][]int{
		{1, 1},
		{0, 1},
		{1, 0},
		{0, 0},
	}

	counter := 0

	for j:=0; j<8; j+= 2 {
		ch := (b >> (7-j)) & 0x01
		ch2 := (b >> (7-(j+1))) & 0x01
		
		if fmt.Sprint(values[counter][0]) != fmt.Sprint(ch) {
			t.Fatal("unexpected value")
		} 
		if fmt.Sprint(values[counter][1]) != fmt.Sprint(ch2) {
			t.Fatal("unexpected value")
		} 
		counter ++
	}
}