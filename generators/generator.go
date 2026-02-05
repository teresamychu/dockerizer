package generators

// Generator produces a Dockerfile for a specific language.

// this is not necessary right now but if i want to support other langugages, will need this.
type Generator interface {
	// Generate creates Dockerfile content for the given project path.
	GenerateDockerfile(projectPath string) (string, error)
	GenerateDockerCompose(projectPath string) (string, error)
}
