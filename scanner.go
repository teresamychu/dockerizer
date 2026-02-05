package main

import (
	"maps"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/teresamychu/dockerizer/generators"
	"golang.org/x/mod/modfile"
	"gopkg.in/yaml.v3"
)

// ImportsMap maps Go package imports to services
// Uses latest available image of each detected service. - this should probably be CLI configurable at some point.
var ImportsMap = map[string]generators.ServiceInfo{
	"lib/pq":       {Name: "postgres", Image: "postgres:16-alpine", Port: "5432"},
	"jackc/pgx":    {Name: "postgres", Image: "postgres:16-alpine", Port: "5432"},
	"mysql":        {Name: "mysql", Image: "mysql:8", Port: "3306"},
	"go-redis":     {Name: "redis", Image: "redis:alpine", Port: "6379"},
	"mongo-driver": {Name: "mongo", Image: "mongo:latest", Port: "27017"},
}

// ParseGoMod extracts module info and detects services from go.mod
func ParseGoMod(path string) (moduleName string, goVersion string, services []generators.ServiceInfo) {
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
		for prefix, svc := range ImportsMap {
			if strings.Contains(require.Mod.Path, prefix) {
				services = append(services, svc)
				break
			}
		}
	}

	return moduleName, goVersion, services
}

// ScanEnvFiles reads .env files and returns key-value pairs
func ScanEnvFiles(projectPath string) map[string]string {
	envVars := make(map[string]string)

	matchFiles, _ := filepath.Glob(filepath.Join(projectPath, ".env*"))

	for _, file := range matchFiles {
		name := filepath.Base(file)
		if name != ".env" && !strings.HasPrefix(name, ".env.") {
			continue
		}
		fileVars, err := godotenv.Read(file)
		if err == nil {
			maps.Copy(envVars, fileVars)
		}
	}

	return envVars
}

// ScanYamlFiles reads yaml/yml files and returns flattened key-value pairs
func ScanYamlFiles(projectPath string) map[string]string {
	configMap := make(map[string]string)

	yamlFiles, _ := filepath.Glob(filepath.Join(projectPath, "*.yaml"))
	ymlFiles, _ := filepath.Glob(filepath.Join(projectPath, "*.yml"))
	allFiles := append(yamlFiles, ymlFiles...)

	for _, f := range allFiles {
		data, err := os.ReadFile(f)
		if err != nil {
			continue
		}
		var yamlMap map[string]interface{}
		if err := yaml.Unmarshal(data, &yamlMap); err != nil {
			continue
		}
		for key, value := range yamlMap {
			switch v := value.(type) {
			case string:
				configMap[key] = v
			case int:
				configMap[key] = strconv.Itoa(v)
			case bool:
				configMap[key] = strconv.FormatBool(v)
			}
		}
	}

	return configMap
}

// ScanConfig scans all config sources and returns AppConfig
func ScanConfig(projectPath string) generators.AppConfig {
	// Scan all sources
	envVars := ScanEnvFiles(projectPath)
	yamlVars := ScanYamlFiles(projectPath)

	// Merge (env takes priority)
	allVars := make(map[string]string)
	maps.Copy(allVars, yamlVars)
	maps.Copy(allVars, envVars)

	// Extract port
	config := generators.AppConfig{
		EnvVars: allVars,
	}

	portKeys := []string{"PORT", "port"}
	for _, key := range portKeys {
		if port, ok := allVars[key]; ok && port != "" {
			config.Port = port
			break
		}
	}

	return config
}
