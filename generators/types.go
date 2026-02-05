package generators

// AppInfo contains project info from go.mod
type AppInfo struct {
	BinaryName string
	GoVersion  string
	HasGoSum   bool
	Services   []ServiceInfo
}

// AppConfig contains config from .env and yaml files
type AppConfig struct {
	Port    string
	EnvVars map[string]string
}

// ServiceInfo represents a detected service dependency.
type ServiceInfo struct {
	Name   string
	Image  string
	Port   string
	EnvVar map[string]string
}

// ServiceDef defines a service in the config file
type ServiceDef struct {
	Imports      []string `yaml:"imports"`
	CodePatterns []string `yaml:"code_patterns"`
	Image        string   `yaml:"image"`
	Port         string   `yaml:"port"`
}
