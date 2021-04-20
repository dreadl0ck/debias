package debias

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
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
		pr  *io.PipeReader
		out string
		ctx context.Context
	)
	if mode == ModeKaminsky {
		pr, ctx, _ = Kaminsky(reader, false, 512)
		out = filepath.Join(path, finfo.Name()+"-kaminsky-debiased.bin")
	} else {
		pr, ctx, _ = VonNeumann(reader, false)
		out = filepath.Join(path, finfo.Name()+"-neumann-debiased.bin")
	}

	// wait until processing is complete
	<-ctx.Done()

	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(pr)
	if err != nil {
		log.Fatal(err)
	}

	// write output buffer
	n, err := f.Write(data)
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
