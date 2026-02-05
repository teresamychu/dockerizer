package generators

import (
	"bytes"
	_ "embed"
	template2 "text/template"
)

//go:embed templates/go.tmpl
var golangTemp string

//go:embed templates/compose.tmpl
var composeTemp string

type GoGenerator struct {
	AppInfo
	AppConfig
}

func NewGoGenerator() *GoGenerator {
	return &GoGenerator{}
}

func (g *GoGenerator) GenerateDockerfile(projectPath string) (string, error) {
	template, _ := template2.New("golang").Parse(golangTemp)

	var buf bytes.Buffer
	err := template.Execute(&buf, g)
	if err != nil {
		return "", err
	}
	return buf.String(), nil

}

func (g *GoGenerator) GenerateComposeFile(projectPath string) (string, error) {
	composeTemplate, _ := template2.New("compose").Parse(composeTemp)

	var compose bytes.Buffer
	err := composeTemplate.Execute(&compose, g)
	if err != nil {
		return "", err
	}
	return compose.String(), nil
}
