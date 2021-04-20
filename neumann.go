package debias

import (
	"bytes"
	"context"
	"io"
)

// The Von Neumann Debiasing algorithm works on pairs of bits, and produces output as follows:
// - If the input is "00" or "11", the input is discarded (no output).
// - If the input is "10", output a "1".
// - If the input is "01", output a "0".
func VonNeumann(reader io.ByteReader) (*bytes.Buffer, context.Context, context.CancelFunc) {
	var (
		buf         bytes.Buffer
		outByte     = byte(0)
		bitCount    uint
		ctx, cancel = context.WithCancel(context.Background())
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				// write leftover
				buf.WriteByte(outByte)
				return
			default:
				b, err := reader.ReadByte()
				if err != nil {

					// write leftover
					buf.WriteByte(outByte)
					cancel()

					return
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

						buf.WriteByte(outByte)
						//fmt.Printf("%08b\n", outByte)
						outByte = byte(0)
					}
				}
			}
		}
	}()

	return &buf, ctx, cancel
}
