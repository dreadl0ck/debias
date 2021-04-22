package debias

import (
	"bufio"
	"fmt"
	"io"
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
		pr  *io.PipeReader
		out string
	)
	if mode == ModeKaminsky {
		pr, _, _ = Kaminsky(reader, false, 512)
		out = filepath.Join(path, finfo.Name()+"-kaminsky-debiased.bin")
	} else {
		pr, _, _ = VonNeumann(reader, false)
		out = filepath.Join(path, finfo.Name()+"-neumann-debiased.bin")
	}

	f, err := os.Create(out)
	if err != nil {
		log.Fatal(err)
	}

	var numBytesWritten int

	for {
		var data = make([]byte, MaxChunkSize)
		n, err := pr.Read(data)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				fmt.Println(err)
				break
			}
			log.Fatal(err)
		}
		data = data[:n]

		// write output buffer
		n, err = f.Write(data)
		if err != nil {
			log.Fatal(err)
		}

		numBytesWritten += n
	}

	dur := time.Since(start)
	fmt.Println("wrote", numBytesWritten, "bytes to output file", out, "in", dur)

	// close output file handle
	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	// return stats to caller
	return &Stats{
		FileName: finfo.Name(),
		BytesIn:  finfo.Size(),
		BytesOut: int64(numBytesWritten),
		Duration: dur,
	}
}
