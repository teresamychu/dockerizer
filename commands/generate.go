package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/teresamychu/dockerizer/generators"
	"github.com/teresamychu/dockerizer/scanner"
)

type GenerateCommand struct{}

func (c *GenerateCommand) Help() string {

	return "Usage: dockerize generate" +
		"Generates dockerfile and docker-compose.yml"

}

func (c *GenerateCommand) Synopsis() string {
	return "Generate Dockerfile and docker-compose.yml"
}

func (c *GenerateCommand) Run(args []string) int {
	flags := flag.NewFlagSet("generate", flag.ContinueOnError)
	dockerfileOutput := flags.String("dockerfile", "Dockerfile", "Dockerfile output filename")
	composeOutput := flags.String("compose", "docker-compose.yml", "Compose output filename")

	if err := flags.Parse(args); err != nil {
		return 1
	}

	//lets assume we're always in the source directory.
	path := "."
	if flags.NArg() > 0 {
		path = flags.Arg(0)
	}

	module, goversion, services := scanner.ParseGoMod(filepath.Join(path, "go.mod"))
	hasGoSum := scanner.FileExists(filepath.Join(path, "go.sum"))

	appConfig := scanner.ScanFiles(path)

	//weve scanned everything and have not found a port. Assume 8080. Probably configurable later?
	if appConfig.Port == "" {
		appConfig.Port = "8080"
	}

	gen := generators.GoGenerator{
		AppInfo: generators.AppInfo{
			BinaryName: module,
			GoVersion:  goversion,
			HasGoSum:   hasGoSum,
			Services:   services,
		},
		AppConfig: appConfig,
	}

	dockerFile, err := gen.GenerateDockerfile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n\n", err)
		return 1
	}
	outputPath := filepath.Join(path, *dockerfileOutput)
	if err := os.WriteFile(outputPath, []byte(dockerFile), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing file: %v\n\n", err)
		return 1
	}
	composeFile, err := gen.GenerateComposeFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating docker-compose: %v\n\n", err)
		return 1
	}
	composeOutputPath := filepath.Join(path, *composeOutput)
	err = os.WriteFile(composeOutputPath, []byte(composeFile), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error writing docker compose file: %v\n\n", err)
		return 1
	}

	fmt.Printf("Generated %s for %s", outputPath, goversion)
	return 0
}
