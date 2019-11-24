package model

import (
	"fmt"
	"io"
)

// Package is a go-package special description.
type Package struct {
	SrcDir string
	Name   string
	Init   string
	Tests  []string
}

func (pkg *Package) Print(w io.Writer) {
	fmt.Fprintf(w, "package %s\n", pkg.Name)
}
