package debias_test

import (
	"fmt"
	"github.com/dreadl0ck/debias"
	"log"
	"os"
	"testing"
	"time"
)

func TestReadFile(t *testing.T) {
	file := "data/data-2200000000.wav"
	fi, err := os.Stat(file)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()
	s := debias.File("data", file, fi, debias.ModeVonNeumann)
	if s == nil {
		t.Fatal(err)
	}
	fmt.Println(time.Since(start))
}
