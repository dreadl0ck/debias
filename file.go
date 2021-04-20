package debias

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

// File will debias a file using the chosen method
func File(path string, file string, finfo fs.FileInfo, mode Mode) *Stats {

	start := time.Now()

	inFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}
	defer inFile.Close()

	reader := bufio.NewReader(inFile)

	var (
		buf *bytes.Buffer
		out string
		ctx context.Context
	)
	if mode == ModeKaminsky {
		buf, ctx, _ = Kaminsky(reader, 512)
		out = filepath.Join(path, finfo.Name()+"-kaminsky-debiased.bin")
	} else {
		buf, ctx, _ = VonNeumann(reader)
		out = filepath.Join(path, finfo.Name()+"-neumann-debiased.bin")
	}

	// wait until processing is complete
	<-ctx.Done()

	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	// write output buffer
	n, err := f.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	dur := time.Since(start)
	fmt.Println("wrote", n, "bytes to output file", out, "in", dur)

	// close output file handle
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// return stats to caller
	return &Stats{
		FileName: finfo.Name(),
		BytesIn:  finfo.Size(),
		BytesOut: int64(n),
		Duration: dur,
	}
}
