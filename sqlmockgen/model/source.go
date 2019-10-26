package model

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	testfile = "_test.go"
	tempfile = testfile + ".tmp.go"
)

func AvoidTesting(path string) (dir string, err error) {
	dir, err = forEachSrcFile(
		path,
		testfile,
		func(filename string) error {
			return copyFile(
				filename,
				filename+tempfile)
		},
	)
	return
}

func RestoreTesting(path string) error {
	_, err := forEachSrcFile(path, tempfile, os.Remove)
	return err
}

func forEachSrcFile(
	path, suffix string,
	do func(filename string) error,
) (
	dir string,
	err error,
) {
	dir, err = sourceDir(path)
	if err != nil {
		return
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, f := range files {
		if f.Mode()&os.ModeType != 0 {
			continue
		}
		if !strings.HasSuffix(f.Name(), suffix) {
			continue
		}

		err = do(filepath.Join(dir, f.Name()))
		if err != nil {
			return
		}
	}

	return
}

func sourceDir(path string) (string, error) {
	var mode build.ImportMode
	pkg, err := build.Default.Import(path, ".", mode)
	if err != nil {
		return "", err
	}
	return pkg.Dir, nil
}

func copyFile(src, dst string) error {
	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dst, data, 0600)
	if err != nil {
		return err
	}

	return nil
}
