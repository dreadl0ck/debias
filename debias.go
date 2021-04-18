package debias

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Mode int

const (
	ModeVonNeumann = 0
	ModeKaminsky = 1
)

// Directory will debias all files in the given directory
// and only for files that have the given extension.
func Directory(path string, ext string, mode Mode) {
	
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	for _, f := range files {
		
		if filepath.Ext(f.Name()) != ext {
			continue
		}
		
		file := filepath.Join(path, f.Name())
		fmt.Println("processing", file)
		
		File(path, file, f, mode)
	}
}

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
func VonNeumann(data []byte) bytes.Buffer {
	var (
		buf bytes.Buffer
		outByte = byte(0)
		bitcount uint
	)

	for _, b := range data {
	
		for j:=0; j<8; j+= 2 {
			
			ch := (b >> (7-j)) & 0x01
			ch2 := (b >> (7-(j+1))) & 0x01
			
			if (ch != ch2) {
			
				if ch == 1 {
					// store a 1 in our bitbuffer
					outByte = setBit(outByte, 7-bitcount)
				} // else: leave the buffer alone, it's already 0 at this bit
				
				bitcount++
			}

			// is the byte full?
			if bitcount == 8 {
				bitcount = 0
				
				buf.WriteByte(outByte)
				outByte = byte(0)
			}
		}
	}

	// write leftover
	buf.WriteByte(outByte)

	return buf
}

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
// 
// Kaminsky:
// - collect discarded bytes
// - use discarded bytes as input for SHA512
// - use the SHA512 hash as key for encrypting the output data with AES
func Kaminsky(data []byte) (bytes.Buffer) {
	var (
		buf bytes.Buffer
		discardBuf bytes.Buffer
		
		// discard byte
		discardByte = byte(0)
		discardBitCount uint

		// out byte
		outByte = byte(0)
		outBitCount uint
	)

	for _, b := range data {
	
		for j:=0; j<8; j+= 2 {
			
			ch := (b >> (7-j)) & 0x01
			ch2 := (b >> (7-(j+1))) & 0x01
						
			if (ch != ch2) {
			
				if ch == 1 {
					// store a 1 in our bitbuffer
					outByte = setBit(outByte, 7-outBitCount)

				} // else: leave the buffer alone, it's already 0 at this bit
				
				outBitCount++
			} else {
				// discarded bits: collect
				discardByte = setBit(discardByte, 7-discardBitCount)
				discardBitCount++
			}

			if discardBitCount == 8 {
				discardBitCount = 0
				
				discardBuf.WriteByte(discardByte)
				discardByte = byte(0)
			}

			// is the byte full?
			if outBitCount == 8 {
				outBitCount = 0
								
				buf.WriteByte(outByte)
				outByte = byte(0)
			}
		}
	}

	// write leftover
	buf.WriteByte(outByte)
	discardBuf.WriteByte(discardByte)

	fmt.Println("kaminsky mode, discard buffer:", discardBuf.Len())

	// create SHA256
	h := sha256.Sum256(discardBuf.Bytes())
	
	// convert into []byte to please go compiler
	var key = make([]byte, 32)
	for i:=0; i<32; i++ {
		key[i] = h[i]
	}

	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}

	// pad plaintext
	bPlaintext := PKCS5Padding(buf.Bytes(), aes.BlockSize, buf.Len())

	// init cipher with key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}
	
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)

	buf.Reset()
	buf.Write(ciphertext)

	return buf
}

func PKCS5Padding(ciphertext []byte, blockSize int, after int) []byte {
	padding := (blockSize - len(ciphertext)%blockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// File will debias a file using the chosen method
func File(path string, file string, finfo fs.FileInfo, mode Mode) {

	start := time.Now()
	
	data, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("read", len(data), "bytes from file", file)

	var buf bytes.Buffer
	if mode == ModeKaminsky {
		buf = Kaminsky(data)
	} else {
		buf = VonNeumann(data)
	}
	
	var out string
	if mode == ModeKaminsky {
		out = filepath.Join(path, finfo.Name() + "-ka-debiased.bin")
	} else {
		out = filepath.Join(path, finfo.Name() + "-vn-debiased.bin")
	}
	
	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	n, err := f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("wrote", n, "bytes to output file", out, "in", time.Since(start))

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func setBit(n byte, pos uint) byte {
    n |= (1 << pos)
    return n
}