package scanner

import (
	_ "embed"
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

//go:embed config/services.yaml
var servicesYAML []byte

// ServiceDefs loaded from services.yaml
var ServiceDefs = LoadServiceDefs()

func LoadServiceDefs() map[string]generators.ServiceDef {
	defs := make(map[string]generators.ServiceDef)
	yaml.Unmarshal(servicesYAML, &defs)
	return defs
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
	services = DetectServicesFromImports(parsedFile)

	return moduleName, goVersion, services
}

// DetectServicesFromImports matches go.mod requires against ServiceDefs
func DetectServicesFromImports(parsedFile *modfile.File) []generators.ServiceInfo {
	var services []generators.ServiceInfo
	seen := make(map[string]bool)

	for _, require := range parsedFile.Require {
		for serviceName, def := range ServiceDefs {
			if seen[serviceName] {
				continue
			}
			for _, pattern := range def.Imports {
				if strings.Contains(require.Mod.Path, pattern) {
					services = append(services, generators.ServiceInfo{
						Name:  serviceName,
						Image: def.Image,
						Port:  def.Port,
					})
					seen[serviceName] = true
					break
				}
			}
		}
	}

	return services
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

// ScanFiles scans all config sources and returns AppConfig
func ScanFiles(projectPath string) generators.AppConfig {
	// Scan all sources
	envVars := ScanEnvFiles(projectPath)
	yamlVars := ScanYamlFiles(projectPath)

	// Merge (env takes priority)
	allVars := make(map[string]string)
	maps.Copy(allVars, yamlVars)
	maps.Copy(allVars, envVars)

	config := generators.AppConfig{
		EnvVars: allVars,
	}

	// Extract port
	//TODO: i feel like this should be kept in the yaml file too.
	portKeys := []string{"PORT", "port"}
	for _, key := range portKeys {
		if port, ok := allVars[key]; ok && port != "" {
			config.Port = port
			break
		}
	}

	return config
}
