package main

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func generate(dsc *Descr) error {
	// TODO: sanity check arguments

	f, err := ioutil.TempFile(dsc.Pkg.SrcDir, "*_test.go")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	entrypoint, err := program(dsc, f)
	if err != nil {
		return err
	}

	return run(entrypoint, dsc.Pkg.SrcDir)
}

func run(entrypoint, dir string) error {
	cmd := exec.Command("go", "test", "-count", "1", "-run", entrypoint, ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
