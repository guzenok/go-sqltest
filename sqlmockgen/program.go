package main

import (
	"bytes"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"text/template"
)

var (
	progOnly   = flag.Bool("prog_only", false, "Only make the generation program; write it to stdout and exit.")
	buildFlags = flag.String("build_flags", "", "Additional flags for go build.")
)

func generate(descr *Descr, destination string) error {
	// TODO: sanity check arguments

	program, err := writeProgram(descr)
	if err != nil {
		return err
	}

	if *progOnly {
		os.Stdout.Write(program)
		os.Exit(0)
	}

	// Try to run the program in the same directory as the input package.
	if err := runInDir(program, descr.Pkg.SrcDir, destination); err == nil {
		return nil
	}
	// Since that didn't work, try to run it in the current working directory.
	wd, _ := os.Getwd()
	if err := runInDir(program, wd, destination); err == nil {
		return nil
	}
	// Since that didn't work, try to run it in a standard temp directory.
	return runInDir(program, "", destination)
}

var generatorSrcCode = template.Must(
	template.New("program").Parse(srcCodeTemplate),
)

func writeProgram(descr *Descr) ([]byte, error) {
	program := new(bytes.Buffer)

	if err := generatorSrcCode.Execute(program, descr); err != nil {
		return nil, err
	}
	return program.Bytes(), nil
}

// runInDir writes the given program into the given dir, and runs it there
// with "-output filename" args.
func runInDir(program []byte, dir, destination string) error {
	// We use TempDir instead of TempFile so we can control the filename.
	tmpDir, err := ioutil.TempDir(dir, "sqlmock_gen_")
	if err != nil {
		return err
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			log.Printf("failed to remove temp directory: %s", err)
		}
	}()
	const progSource = "prog.go"
	var progBinary = "prog.bin"
	if runtime.GOOS == "windows" {
		// Windows won't execute a program unless it has a ".exe" suffix.
		progBinary += ".exe"
	}

	if err := ioutil.WriteFile(filepath.Join(tmpDir, progSource), program, 0600); err != nil {
		return err
	}

	cmdArgs := []string{}
	cmdArgs = append(cmdArgs, "build")
	if *buildFlags != "" {
		cmdArgs = append(cmdArgs, *buildFlags)
	}
	cmdArgs = append(cmdArgs, "-o", progBinary, progSource)

	// Build the program.
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = tmpDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return run(filepath.Join(tmpDir, progBinary), destination)
}

// run the given program with "-output filename" args.
func run(program, destination string) error {
	cmd := exec.Command(program, "-output", destination)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
