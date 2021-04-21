package debias_test

import (
	"github.com/dreadl0ck/debias"
	"log"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	file := "data/data-2200000000.wav"
	fi, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	s := debias.File("data", file, fi, debias.ModeVonNeumann)
	if s == nil {
		t.Fatal(err)
	}
}
