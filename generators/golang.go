package generators

import (
	"bytes"
	_ "embed"
	template2 "text/template"
)

//go:embed templates/go.tmpl
var golangTemp string

type GoGenerator struct {
	BinaryName string
	GoVersion  string
	HasGoSum   bool
}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

func (g *GoGenerator) Generate(projectPath string) (string, error) {
	// TODO: Implement using templates/go.tmpl
	template, _ := template2.New("golang").Parse(golangTemp)

	var buf bytes.Buffer
	err := template.Execute(&buf, g)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
