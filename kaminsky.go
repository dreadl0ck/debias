package debias

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
)

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
//
// Kaminsky addition:
// - collect discarded bytes
// - use discarded bytes as input for SHA512
// - use the SHA512 hash as key for encrypting the output data with AES
func Kaminsky(reader io.ByteReader, wait bool, blockSize int64) (*io.PipeReader, context.Context, context.CancelFunc) {
	var (
		// internal buffer
		buf        bytes.Buffer
		discardBuf bytes.Buffer

		// discard byte
		discardByte     = byte(0)
		discardBitCount uint

		// out byte
		outByte     = byte(0)
		outBitCount uint

		ctx, cancel = context.WithCancel(context.Background())

		numBytes int64
		pr, pw   = io.Pipe()
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				// write leftover
				buf.WriteByte(outByte)
				discardBuf.WriteByte(discardByte)
				aesEncrypt(&buf, &discardBuf, pw)
				err := pw.Close()
				if err != nil {
					log.Println("failed to close pipe writer:", err)
				}
				return
			default:
				b, err := reader.ReadByte()
				if err != nil {

					if !wait {
						// write leftover
						buf.WriteByte(outByte)
						discardBuf.WriteByte(discardByte)
						aesEncrypt(&buf, &discardBuf, pw)
						err := pw.Close()
						if err != nil {
							log.Println("failed to close pipe writer:", err)
						}
						cancel()
						return
					}
				}

				numBytes++
				if numBytes%blockSize == 0 {
					aesEncrypt(&buf, &discardBuf, pw)
				}

				for j := 0; j < 8; j += 2 {

					ch := (b >> (7 - j)) & 0x01
					ch2 := (b >> (7 - (j + 1))) & 0x01

					if ch != ch2 {

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
		}
	}()

	return pr, ctx, cancel
}

// this function:
// - encrypts the input buffer with aes-256-cbc
// - reads the IV via crypto/rand from /dev/random for every invocation
// - calculates the key based on the key buffer
// - writes the result in the output buffer
// - resets the input buffer, but never the key buffer
func aesEncrypt(buf *bytes.Buffer, keyBuf *bytes.Buffer, pw *io.PipeWriter) {

	// create SHA256
	h := sha256.Sum256(keyBuf.Bytes())

	// convert into []byte to please go compiler
	var key = make([]byte, 32)
	for i := 0; i < 32; i++ {
		key[i] = h[i]
	}

	// read random bytes from crypto/rand for the initialization vector
	iv := make([]byte, 16)
	_, err := rand.Read(iv)
	if err != nil {
		log.Fatal(err)
	}

	// pad plaintext to aes.BlockSize
	bPlaintext := pkcs5Padding(buf.Bytes(), aes.BlockSize)

	// init AES cipher with SHA256 as key
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Fatal(err)
	}

	// create CBC block mode and encrypt data
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)

	// reset internal buffer
	buf.Reset()

	// never reset key buffer, so it keeps growing

	// write output
	pw.Write(ciphertext)
}
