package generators

// Generator produces a Dockerfile for a specific language.

type Generator interface {
	// Generate creates Dockerfile content for the given project path.
	Generate(projectPath string) (string, error)
}
