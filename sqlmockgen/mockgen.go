package main

import (
	"flag"
	"go/build"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	gomockImportPath = "github.com/guzenok/go-sqltest/gomock"
)

var (
	destination     = flag.String("destination", "", "Output file; defaults to stdout.")
	mockNames       = flag.String("mock_names", "", "Comma-separated interfaceName=mockName pairs of explicit mock names to use. Mock names default to 'Mock'+ interfaceName suffix.")
	packageOut      = flag.String("package", "", "Package of the generated code; defaults to the package of the input with a 'mock_' prefix.")
	selfPackage     = flag.String("self_package", "", "The full package import path for the generated code. The purpose of this flag is to prevent import cycles in the generated code by trying to include its own package. This can happen if the mock's package is set to one of its inputs (usually the main one) and the output is stdio so sqlmockgen cannot detect the final output package. Setting this flag will then tell sqlmockgen which import to exclude.")
	writePkgComment = flag.Bool("write_package_comment", true, "Writes package documentation comment (godoc) if true.")
	copyrightFile   = flag.String("copyright_file", "", "Copyright file used to add copyright header")

	debugParser = flag.Bool("debug_parser", false, "Print out parser results only.")
)

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() != 2 {
		usage()
		log.Fatal("Expected exactly two arguments")
	}
	pkg, err := reflect(flag.Arg(0), strings.Split(flag.Arg(1), ","))
	if err != nil {
		log.Fatalf("Loading input failed: %v", err)
	}

	if *debugParser {
		pkg.Print(os.Stdout)
		return
	}

	dst := os.Stdout
	if len(*destination) > 0 {
		if err := os.MkdirAll(filepath.Dir(*destination), os.ModePerm); err != nil {
			log.Fatalf("Unable to create directory: %v", err)
		}
		f, err := os.Create(*destination)
		if err != nil {
			log.Fatalf("Failed opening destination file: %v", err)
		}
		defer f.Close()
		dst = f
	}

	packageName := *packageOut
	if packageName == "" {
		// pkg.Name in reflect mode is the base name of the import path,
		// which might have characters that are illegal to have in package names.
		packageName = "mock_" + sanitize(pkg.Name)
	}

	// outputPackagePath represents the fully qualified name of the package of
	// the generated code. Its purposes are to prevent the module from importing
	// itself and to prevent qualifying type names that come from its own
	// package (i.e. if there is a type called X then we want to print "X" not
	// "package.X" since "package" is this package). This can happen if the mock
	// is output into an already existing package.
	outputPackagePath := *selfPackage
	if len(outputPackagePath) == 0 && len(*destination) > 0 {
		dst, _ := filepath.Abs(filepath.Dir(*destination))
		for _, prefix := range build.Default.SrcDirs() {
			if strings.HasPrefix(dst, prefix) {
				if rel, err := filepath.Rel(prefix, dst); err == nil {
					outputPackagePath = rel
					break
				}
			}
		}
	}

	g := new(generator)
	g.srcPackage = flag.Arg(0)
	g.srcInterfaces = flag.Arg(1)

	if *mockNames != "" {
		g.mockNames = parseMockNames(*mockNames)
	}
	if *copyrightFile != "" {
		header, err := ioutil.ReadFile(*copyrightFile)
		if err != nil {
			log.Fatalf("Failed reading copyright file: %v", err)
		}

		g.copyrightHeader = string(header)
	}
	if err := g.Generate(pkg, packageName, outputPackagePath); err != nil {
		log.Fatalf("Failed generating mock: %v", err)
	}
	if _, err := dst.Write(g.Output()); err != nil {
		log.Fatalf("Failed writing to destination: %v", err)
	}
}

func parseMockNames(names string) map[string]string {
	mocksMap := make(map[string]string)
	for _, kv := range strings.Split(names, ",") {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 || parts[1] == "" {
			log.Fatalf("bad mock names spec: %v", kv)
		}
		mocksMap[parts[0]] = parts[1]
	}
	return mocksMap
}

func usage() {
	io.WriteString(os.Stderr, usageText)
	flag.PrintDefaults()
}

const usageText = `sqlmockgen generates mock interfaces by building a program
that uses reflection to understand interfaces. It is enabled
by passing two non-flag arguments: an import path, and a
comma-separated list of symbols.
Example:
	sqlmockgen database/sql/driver Conn,Driver

`
