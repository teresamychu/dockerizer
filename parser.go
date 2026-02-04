package main

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

func parseGoMod(path string) (moduleName string, goVersion string, services []ServiceInfo) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "app", "1.22", nil
	}
	parsedFile, err := modfile.Parse("go.mod", data, nil)
	if err != nil {
		return "app", "1.22", nil
	}

	moduleName = filepath.Base(parsedFile.Module.Mod.Path)
	goVersion = parsedFile.Go.Version

	for _, require := range parsedFile.Require {
		for prefix, svc := range ServiceMap {
			if strings.Contains(require.Mod.Path, prefix) {
				services = append(services, svc)
				break
			}
		}
		if svc, ok := ServiceMap[require.Mod.Path]; ok {
			services = append(services, svc)
		}
	}

	return moduleName, goVersion, services

}
