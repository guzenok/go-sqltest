package generator

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func NewFile(filename string) (w io.Writer, close func(), err error) {
	w = os.Stdout
	close = func() {}

	if len(filename) == 0 {
		return
	}

	err = os.MkdirAll(filepath.Dir(filename), os.ModePerm)
	if err != nil {
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		return
	}
	w = f
	close = func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close output file %q", filename)
		}
	}

	return
}
