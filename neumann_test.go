package debias_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/dreadl0ck/debias"
)

func TestNeumannEntropy(t *testing.T) {
	var b = []byte{179, 14, 102, 35, 184, 17, 183, 66, 242, 171, 212, 0, 202, 9, 185, 88, 146, 157, 132, 67, 81, 180, 30, 236, 205, 236, 233, 14, 180, 75, 209, 224, 101, 207, 237, 44, 68, 54, 196, 188, 129, 23, 240, 254, 14, 2, 32, 105, 223, 237, 1, 145, 243, 122, 57, 93, 213, 178, 10, 197, 186, 172, 103, 185, 86, 30, 117, 184, 37, 163, 56, 153, 89, 119, 153, 66, 107, 185, 98, 141, 92, 175, 31, 174, 90, 172, 252, 172, 191, 111, 38, 142, 40, 227, 146, 195, 67, 114, 89, 158, 1, 226, 248, 145, 121, 223, 32, 6, 226, 4, 93, 3, 5, 101, 239, 159, 66, 53, 71, 4, 3, 212, 233, 59, 236, 193, 4, 182, 183, 183, 68, 127, 134, 69, 198, 221, 41, 25, 114, 85, 32, 188, 195, 154, 224, 222, 62, 11, 171, 70, 44, 82, 6, 65, 184, 144, 178, 139, 6, 109, 136, 69, 195, 1, 215, 238, 242, 235, 96, 46, 101, 134, 2, 133, 222, 246, 36, 95, 53, 32, 241, 254, 219, 52, 11, 43, 6, 42, 197, 191, 223, 183, 132, 117, 184, 247, 136, 186, 15, 164, 0, 184, 221, 151, 130, 104, 213, 53, 32, 160, 11, 73, 33, 208, 92, 75, 175, 108, 122, 131, 173, 188, 200, 61, 208, 73, 34, 209, 102, 118, 70, 216, 23, 68, 246, 52, 33, 213, 116, 47, 41, 61, 38, 222, 132, 86, 203, 94, 254, 79, 124, 207, 46, 145, 44, 131, 112, 1, 184, 159, 235, 101, 23, 190, 133, 237, 16, 239, 208, 57, 2, 226, 126, 150, 162, 77, 118, 13, 32, 52, 121, 66, 23, 18, 246, 208, 176, 19, 176, 166, 208, 71, 175, 144, 93, 179, 64, 46, 193, 3, 194, 211, 86, 209, 106, 49, 159, 202, 30, 146, 19, 133, 233, 68, 59, 141, 149, 1, 27, 24, 190, 59, 229, 121, 133, 85, 216, 67, 143, 207, 118, 238, 199, 46, 150, 200, 84, 196, 213, 35, 92, 255, 226, 30, 108, 96, 224, 94, 204, 46, 134, 206, 107, 230, 41, 48, 208, 80, 34, 43, 179, 25, 17, 59, 52, 53, 189, 210, 142, 16, 158, 241, 65, 223, 81, 52, 92, 209, 218, 57, 53, 219, 66, 33, 153, 49, 209, 173, 231, 210, 74, 82, 16, 70, 240, 220, 31, 221, 105, 21, 136, 182, 239, 165, 147, 104, 102, 53, 57, 137, 181, 197, 205, 26, 126, 54, 50, 60, 228, 83, 98, 213, 61, 172, 95, 131, 83, 199, 10, 193, 164, 115, 39, 89, 71, 71, 76, 95, 65, 32, 43, 128, 94, 117, 46, 177, 174, 155, 137, 77, 144, 84, 9, 46, 249, 33, 174, 50, 37, 25, 79, 29, 247, 30, 35, 229, 98, 34, 157, 235, 125, 10, 36, 53, 254, 154, 90, 228, 45, 175, 99, 185, 112, 153, 205, 158, 163, 31, 81, 69, 234, 225, 184, 0, 245, 13, 190, 64, 150, 69, 133, 6, 26, 220, 110, 8, 177, 115, 253, 56, 229, 38, 146, 149, 45, 56, 9, 212, 33, 31, 142, 78, 133, 10, 37, 61, 154, 52, 131, 1, 117, 202, 148, 84, 52, 15, 177, 39, 97, 242, 199, 83, 151, 194, 21, 166, 114, 98, 72, 215, 131, 52, 90, 126, 87, 165, 184, 83, 7, 236, 107, 78, 143, 214, 151, 157, 84, 118, 20, 186, 33, 178, 179, 18, 172, 190, 102, 242, 247, 94, 97, 58, 102, 72, 236, 147, 1, 141, 69, 117, 135, 227, 36, 242, 251, 214, 224, 135, 168, 184, 201, 42, 46, 95, 137, 167, 15, 48, 199, 96, 108, 145, 16, 40, 55, 251, 246, 4, 15, 16, 123, 41, 78, 110, 69, 60, 253, 251, 218, 26, 78, 17, 44, 210, 46, 30, 7, 50, 130, 66, 60, 249, 49, 177, 69, 254, 219, 214, 89, 33, 240, 63, 105, 48, 46, 178, 216, 171, 44, 43, 29, 19, 239, 199, 48, 214, 57, 1, 244, 89, 140, 129, 149, 228, 87, 43, 193, 166, 231, 19, 31, 68, 254, 97, 207, 38, 160, 5, 186, 21, 103, 211, 38, 41, 142, 179, 177, 72, 220, 97, 144, 197, 72, 102, 47, 69, 230, 234, 33, 1, 14, 73, 217, 161, 103, 74, 174, 166, 113, 93, 38, 13, 228, 19, 201, 111, 107, 55, 115, 4, 138, 150, 98, 185, 33, 16, 187, 128, 16, 124, 196, 177, 144, 144, 110, 210, 236, 99, 159, 212, 208, 150, 251, 25, 239, 245, 173, 89, 76, 193, 98, 111, 12, 10, 118, 63, 209, 187, 14, 90, 150, 41, 199, 175, 31, 47, 48, 80, 89, 39, 112, 6, 81, 231, 217, 121, 234, 18, 40, 82, 185, 179, 2, 100, 22, 197, 42, 1, 28, 154, 42, 34, 105, 108, 214, 252, 218, 72, 110, 235, 55, 54, 24, 56, 159, 165, 255, 34, 169, 179, 131, 50, 8, 144, 78, 232, 131, 151, 143, 215, 23, 160, 48, 249, 68, 39, 197, 191, 143, 248, 67, 165, 68, 209, 71, 179, 12, 220, 103, 41, 232, 104, 251, 220, 5, 4, 162, 26, 87, 165, 50, 209, 234, 222, 81, 145, 111, 72, 41, 50, 120, 2, 81, 120, 121, 239, 144, 6, 50, 254, 88, 57, 73, 161, 65, 198, 96, 227, 215, 142, 27, 125, 151, 180, 127, 214, 223, 129, 238, 2, 121, 253, 11, 37, 229, 227, 157, 94, 153, 66, 178, 70, 206, 90, 96, 72, 175, 188, 232, 7, 59, 86, 80, 157, 152, 58, 152, 88, 92, 166, 233, 221, 3, 102, 60, 90, 85, 84, 25, 97, 158, 82, 187, 125, 24, 50, 179, 159, 22, 187, 110, 234, 108, 23, 185, 77, 155, 77, 111, 124, 180, 67, 165, 233, 159, 43, 252, 82, 124, 62, 236, 49, 78, 246, 193, 183, 122, 138, 246, 38, 14, 143, 171, 197, 99, 179, 182, 0, 139, 125, 204, 80, 49, 120, 73, 171, 11, 134, 166, 193, 167, 39, 231, 124, 142, 113, 56, 14, 143, 63, 163, 133, 12, 169, 176, 180, 189, 240, 246, 185, 180, 111, 140, 128}
	fmt.Println(len(b), "bytes")

	pr, _, _ := debias.VonNeumann(bytes.NewReader(b), false)

	var data = make([]byte, 300)
	n, err := pr.Read(data)
	if err != nil {
		t.Fatal(err)
	}
	data = data[:n]

	if len(data) != 258 {
		t.Fatal("unexpected number of output bytes, want 258 but got ", len(data))
	}

	if debias.ShannonEntropy(data) < 2000 {
		t.Fatal("entropy too low: got ", debias.ShannonEntropy(data), " expected > ", 2000)
	}
}

// IMPORTANT: both input and output bit strings have to be a multiple of eight!
// An exception to this is an empty string for the output.
type test struct {
	in  string
	out string
}

var tests = []test{
	{
		// 0 to 0000000000000000000 (steps of one) --> nothing
		in:  "0000000000000000000",
		out: "",
	},
	{
		// 1 to 1111111111111111111 (steps of one) --> nothing
		in:  "1111111111111111111",
		out: "",
	},
	{
		// 01 to 01010101 01010101 (steps of two) --> all zeros * 1/2 input length
		in:  "0101010101010101",
		out: "00000000",
	},
	{
		// 10 to 10101010 10101010 (steps of two) --> all ones * 1/2 input length
		in:  "1010101010101010",
		out: "11111111",
	},
	{
		in:  "01101010110011111101100010100110",
		out: "0111011101000000", // contains padding with 6 trailing zeroes because Go's buffers can only hold full bytes.
	},
	{
		in:  "01001011001100100001011010100001",
		out: "0110011100000000", // seven trailing zeroes for padding
	},
	{
		in:  "00111110101000110000100001100011",
		out: "11110100", // two trailing zeroes for padding
	},

	{
		in:  "00111110101000110000100001100011",
		out: "11110100", // two trailing zeroes for padding
	},
	{
		in:  "01001000100100110001101001101101",
		out: "0110011010000000", // six trailing zeroes for padding
	},
	{
		in:  "00111111110011010011100011110011",
		out: "01000000", // six trailing zeroes for padding
	},
	{
		in:  "10110101011001111110110001010011",
		out: "1000101000000000", // six trailing zeroes for padding
	},
	{
		in:  "01100101100110010000101101010000",
		out: "0100101010000000", // five trailing zeroes for padding
	},
	{
		in:  "10011111010100011000010000110001",
		out: "10000100",
	},
	{
		in:  "01100100010010011000110100110110",
		out: "0100101001000000", // six trailing zeroes for padding
	},
	{
		in:  "10011111111001101001110001111001",
		out: "1010110010000000", // six trailing zeroes for padding
	},
	{
		in: "1001100110011001",
		out: "10101010",
	},
}

func TestVonNeumann(t *testing.T) {
	for i, te := range tests {

		// apply von neumann debiasing to input buffer
		pr, _, _ := debias.VonNeumann(bytes.NewBuffer(bitString(te.in).bytes()), false)

		var data = make([]byte, 300)
		n, err := pr.Read(data)
		if err != nil {
			t.Fatal(err)
		}
		data = data[:n]

		// print buffer
		//log.Println(buf.bytes())

		if !bytes.HasPrefix(data, bitString(te.out).bytes()) {

			fmt.Printf("buf.bytes(): %08b\n", data)
			fmt.Printf("bitString(te.out).AsByteSlice(): %08b\n", bitString(te.out).bytes())

			spew.Dump(tests[i])
			fmt.Println("input size", strconv.Itoa(len(te.in))+" bits")
			fmt.Println("output size", strconv.Itoa(len(te.out))+" bits")
			fmt.Println("======================= DEBUG ========================")
			debias.VonNeumann(bytes.NewBuffer(bitString(te.in).bytes()), false)
			fmt.Println("======================= DEBUG ========================")
			t.Fatal("test #", i, ": expected "+te.out+" ("+strconv.Itoa(len(te.out))+" bits) but got "+fmt.Sprintf("%08b", data)+" ("+strconv.Itoa(len(data)*8)+" bits)")
		}
	}
}

func TestPattern(t *testing.T) {

	var (
		b    = bitString("10011001").bytes()
		size = 1000*1000
		data = make([]byte, size)
	)
	for i:=0;i<size;i++ {
		data[i] = b[0]
	}

	// apply von neumann debiasing to input buffer
	pr, _, _ := debias.VonNeumann(bytes.NewBuffer(data), false)

	out, err := ioutil.ReadAll(pr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(out), "bytes", out[:10])
}

func TestDiscardAll(t *testing.T) {

	var (
		b    = bitString("11111111").bytes()
		size = 1000*1000
		data = make([]byte, size)
	)
	for i:=0;i<size;i++ {
		data[i] = b[0]
	}

	// apply von neumann debiasing to input buffer
	pr, _, _ := debias.VonNeumann(bytes.NewBuffer(data), false)

	out, err := ioutil.ReadAll(pr)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(out, "bytes", out[:10])
}

func TestMakeBitPattern1(t *testing.T) {

	f, err := os.Create("10011001.bin")
	if err != nil {
		t.Fatal(err)
	}

	var (
		b    = bitString("10011001").bytes()
		size = 1000*1000*1000
		data = make([]byte, size)
	)
	for i:=0;i<size;i++ {
		data[i] = b[0]
	}

	_, err = f.Write(data)
	if err != nil {
		t.Fatal(err)
	}
}

func TestMakeBitPattern2(t *testing.T) {

	f, err := os.Create("01010101.bin")
	if err != nil {
		t.Fatal(err)
	}

	var (
		b    = bitString("01010101").bytes()
		size = 1000*1000*1000
		data = make([]byte, size)
	)
	for i:=0;i<size;i++ {
		data[i] = b[0]
	}

	_, err = f.Write(data)
	if err != nil {
		t.Fatal(err)
	}
}