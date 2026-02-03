package main

import (
	"os"
	"path/filepath"

	"golang.org/x/mod/modfile"
)

func parseGoMod(path string) (moduleName string, goVersion string) {
	data, _ := os.ReadFile(path)
	parsedFile, _ := modfile.Parse("go.mod", data, nil)

	moduleName = filepath.Base(parsedFile.Module.Mod.Path)
	goVersion = parsedFile.Go.Version

	return moduleName, goVersion

}
