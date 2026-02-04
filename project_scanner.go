package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type AppInfo struct {
	Port     string
	EnvVars  []string
	Services []ServiceInfo
}

type ServiceInfo struct {
	Name  string
	Image string
}

// TODO: make this more versatile and robust.
// Should this map be built dynamically off some image storage?
// Do i have to worry about versions? maybe later.
var ServiceMap = map[string]ServiceInfo{
	"lib/pq":   {Name: "postgres", Image: "postgres:16-alpine"},
	"mysql":    {Name: "mysql", Image: "mysql:8-alpine"},
	"go-redis": {Name: "redis", Image: "redis:8-alpine"},
}

func Scan(projectPath string) (*AppInfo, error) {
	appInfo := &AppInfo{}

	//get all golang files only.
	golangFiles, err := getProjectFiles(projectPath)
	if err != nil {
		return nil, err
	}
	for _, file := range golangFiles {
		if err := scanFile(file, *appInfo); err != nil {
			//dont care if theres an error scanning a single file, keep scanning.
			continue
		}
	}
	return appInfo, nil
}

func getProjectFiles(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, fileInfo os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fileInfo.IsDir() && strings.HasSuffix(fileInfo.Name(), ".go") && !strings.HasSuffix(fileInfo.Name(), "_test.go") {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

func scanFile(filePath string, info AppInfo) error {
	fileSet := token.NewFileSet()
	fileNode, err := parser.ParseFile(fileSet, filePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	//traverse to find env vars and ports? maybe??
	ast.Inspect(fileNode, func(node ast.Node) bool {
		//TODO: use ListenAndServe call to find port - 4 or 5 digits. that sounds like a problem for future teresa
		//TODO: use os.envvars for env vars. BUT HOWWW TERESA?
		return true
	})
	return nil
}
