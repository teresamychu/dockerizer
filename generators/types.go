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

// ServiceInfo represents a detected service dependency
type ServiceInfo struct {
	Name   string
	Image  string
	Port   string
	EnvVar map[string]string
}
