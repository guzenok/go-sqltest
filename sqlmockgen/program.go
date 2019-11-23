package main

import (
	"io/ioutil"
	"os"
	"os/exec"
	"text/template"
)

func generate(dsc *Descr) error {
	// TODO: sanity check arguments

	f, err := ioutil.TempFile(dsc.Pkg.SrcDir, "*_test.go")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())

	specf, err := writeProgram(dsc, f)
	if err != nil {
		return err
	}

	return runInDir(specf, dsc.Pkg.SrcDir)
}

var generatorSrcCode = template.Must(
	template.New("program").Parse(srcCodeTemplate),
)

func writeProgram(dsc *Descr, f *os.File) (funcName string, err error) {
	return dsc.SpecTestName, generatorSrcCode.Execute(f, dsc)
}

func runInDir(specf, dir string) error {
	cmd := exec.Command("go", "test", "-count", "1", "-run", specf, ".")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
