package generator

import (
	"io"
	"log"
	"os"
	"path/filepath"
)

func OutWriter(filename string) (w io.Writer, close func()) {
	w = os.Stdout
	close = func() {}

	if len(filename) == 0 {
		return
	}

	if err := os.MkdirAll(filepath.Dir(filename), os.ModePerm); err != nil {
		log.Fatalf("Unable to create directory: %v", err)
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatalf("failed to open output file %q", filename)
	}

	close = func() {
		if err := f.Close(); err != nil {
			log.Fatalf("failed to close output file %q", filename)
		}
	}

	w = f
	return
}
