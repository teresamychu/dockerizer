package scanner

import (
	"os"
	"path/filepath"
)

type Language string

const (
	LangGo      Language = "go"
	LangNode    Language = "node"
	LangPython  Language = "python"
	LangUnknown Language = "unknown"
)

// Detect identifies the project language based on marker files.
func Detect(projectPath string) Language {
	// Check for: go.mod, package.json, requirements.txt, pyproject.toml
	_ = projectPath
	if fileExists(markerPath(projectPath, "go.mod")) {
		return LangGo
	}

	return LangUnknown
}

func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func fileExists(path string) bool {
	return FileExists(path)
}

func markerPath(projectPath, filename string) string {
	return filepath.Join(projectPath, filename)
}
