package debias

import (
	"bytes"
	"fmt"
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

	dur := time.Since(start)
	fmt.Println("wrote", n, "bytes to output file", out, "in", dur)

	err = f.Close()
	if err != nil {
		log.Fatal(err)
	}

	return &Stats{
		FileName: finfo.Name(),
		BytesIn:  len(data),
		BytesOut: n,
		Duration: dur,
	}
}