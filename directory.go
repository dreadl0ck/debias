package debias

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
)

// Directory will debias all files in the given directory
// and only for files that have the given extension.
func Directory(path string, ext string, mode Mode) []*Stats {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	var stats []*Stats
	for _, f := range files {

		if filepath.Ext(f.Name()) != ext {
			continue
		}

		file := filepath.Join(path, f.Name())
		fmt.Println("processing", file)

		stats = append(stats, File(path, file, f, mode))
	}

	return stats
}