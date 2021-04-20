package debias

import (
	"bytes"
	"context"
	"io"
)

var MaxChunkSize = 1024

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
func VonNeumann(reader io.ByteReader, wait bool) (*io.PipeReader, context.Context, context.CancelFunc) {
	var (
		outByte     = byte(0)
		bitCount    uint
		ctx, cancel = context.WithCancel(context.Background())
		pr, pw      = io.Pipe()

		outBuf bytes.Buffer
	)

	go func() {
		for {
			select {
			case <-ctx.Done():

				// write leftover
				outBuf.WriteByte(outByte)
				pw.Write(outBuf.Bytes())

				return
			default:
				b, err := reader.ReadByte()
				if err != nil {

					if !wait {
						// write leftover
						outBuf.WriteByte(outByte)
						pw.Write(outBuf.Bytes())

						cancel()
						return
					}
				}

				if outBuf.Len() + 1 > MaxChunkSize {
					pw.Write(outBuf.Bytes())
					outBuf.Reset()
				}

				for j := 0; j < 8; j += 2 {

					ch := (b >> (7 - j)) & 0x01
					ch2 := (b >> (7 - (j + 1))) & 0x01

					if ch != ch2 {

						if ch == 1 {
							// store a 1 in our bitbuffer
							outByte = setBit(outByte, 7-bitCount)
						} // else: leave the buffer alone, it's already 0 at this bit

						bitCount++
					}

					// is the byte full?
					if bitCount == 8 {
						bitCount = 0

						outBuf.WriteByte(outByte)

						//fmt.Printf("%08b\n", outByte)
						outByte = byte(0)
					}
				}
			}
		}
	}()

	return pr, ctx, cancel
}
