// Dockerize CLI - generates Dockerfiles for Go, Node.js, and Python projects
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/teresamychu/dockerize/generators"
)

var golangTemplate string

func main() {
	output := flag.String("output", "Dockerfile", "output filename")
	flag.Parse()

	path := "."
	if flag.NArg() > 0 {
		path = flag.Arg(0)
	}

	language := Detect(path)
	if language != LangGo {
		fmt.Println("Not a Golang application")
		os.Exit(1)
	}

	module, goversion := parseGoMod(filepath.Join(path, "go.mod"))

	hasGoSum := fileExists(filepath.Join(path, "go.sum"))

	gen := generators.GoGenerator{
		BinaryName: module,
		GoVersion:  goversion,
		HasGoSum:   hasGoSum,
	}
	dockerFile, err := gen.Generate(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	outputPath := filepath.Join(path, *output)
	err = os.WriteFile(outputPath, []byte(dockerFile), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Generated %s for %s (Go %s)\n", outputPath, module, goversion)
}
