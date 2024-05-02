package main

import (
	_ "embed"
	"fmt"
	"github.com/8naps/power-plugins/plugins/java-springbootapp/generator"
	"github.com/hashicorp/go-plugin"
	"github.com/zcubbs/blueprint"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Generator struct {
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: blueprint.HandshakeConfig,
		Plugins: map[string]plugin.Plugin{
			"blueprint": &blueprint.GeneratorPlugin{Impl: &Generator{}},
		},
	})

	// Hang the main process as the plugin should be run in a separate process
	select {}
}

//go:embed spec.yaml
var specYaml []byte

func (g *Generator) LoadSpec() (blueprint.Spec, error) {
	return blueprint.LoadBlueprintSpecFromBytes(specYaml)
}

// Generate implements the ComponentGenerator interface, generating a Go API server
func (g *Generator) Generate(spec blueprint.Spec, values map[string]string, workdir string) error {
	// Step 1: Parse ComponentSpec to get required configurations
	config, err := parseConfig(spec, values)
	if err != nil {
		return fmt.Errorf("error parsing component spec: %v", err)
	}

	// Step 2: Create project structure and files based on the parsed config
	projectPath, err := createProjectStructure(workdir, config)
	if err != nil {
		return fmt.Errorf("error creating project structure: %v", err)
	}

	// Step 3: Generate project files
	err = generateProjectFiles(projectPath, config)
	if err != nil {
		return fmt.Errorf("error generating project files: %v", err)
	}

	return nil
}

// parseConfig extracts configuration options from ComponentSpec
func parseConfig(spec blueprint.Spec, values map[string]string) (Config, error) {
	var config Config

	// Helper function to set config value with default if needed
	setConfigValue := func(field *string, option blueprint.Option) {
		value := values[option.ID]
		if value == "" {
			*field = option.Default
		} else {
			*field = value
		}
	}

	// Iterate over options and set config values
	for _, option := range spec.Options {
		switch option.ID {
		case "java_version":
			setConfigValue(&config.JavaVersion, option)
		case "packaging":
			setConfigValue(&config.Packaging, option)
		case "springBootVersion":
			setConfigValue(&config.SpringBootVersion, option)
		case "dependencyManagement":
			setConfigValue(&config.DependencyManagement, option)
		case "metadataPackageName":
			setConfigValue(&config.Metadata.PackageName, option)
		case "metadataGroupId":
			setConfigValue(&config.Metadata.GroupId, option)
		case "metadataArtifactId":
			setConfigValue(&config.Metadata.ArtifactId, option)
		case "metadataName":
			setConfigValue(&config.Metadata.Name, option)
		case "metadataDescription":
			setConfigValue(&config.Metadata.Description, option)
		}
	}

	return config, nil
}

// createProjectStructure sets up the project directory and base files
func createProjectStructure(outputPath string, config Config) (string, error) {
	projectPath := filepath.Join(outputPath, "java-springbootapp")
	if err := os.MkdirAll(projectPath, 0750); err != nil {
		return "", err
	}

	// Create subdirectories
	dirs := []string{"src", "src/main", "src/main/java", "src/main/resources", "src/test", "src/test/java"}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0750); err != nil {
			return "", err
		}
	}

	return projectPath, nil
}

// generateProjectFiles generates project files based on the parsed config
func generateProjectFiles(projectPath string, config Config) error {
	// Define the file paths and corresponding templates
	files := map[string]string{
		"pom.xml": generator.PomXmlTemplate,
		"src/main/java/{{.Metadata.PackageName}}/Application.java":     generator.MainJavaTemplate,
		"src/main/java/{{.Metadata.PackageName}}/HelloController.java": generator.HelloControllerJavaTemplate,
	}

	// Process each template and create files
	for filePath, tmpl := range files {
		// Replace placeholders in the file path
		fullPath, err := replacePlaceholdersInString(filePath, config)
		if err != nil {
			return err
		}
		fullPath = filepath.Join(projectPath, filePath)
		if err := processTemplate(fullPath, tmpl, config); err != nil {
			return err
		}
	}

	return nil
}

func processTemplate(fPath, tmpl string, config Config) error {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Clean(fPath))
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	return t.Execute(file, config)
}

// Config contains the configuration options for the generator
type Config struct {
	JavaVersion          string
	Packaging            string
	SpringBootVersion    string
	DependencyManagement string

	Metadata Metadata
}

// Metadata contains the metadata for the generator
type Metadata struct {
	PackageName string
	GroupId     string
	ArtifactId  string
	Description string
	Name        string
}

type StringReplacer string

const (
	CamelCase     StringReplacer = "camelCase"
	DotsToSlashes StringReplacer = "dotsToSlashes"
)

// replacePlaceholdersInString replaces placeholders in a string with values from the config
func replacePlaceholdersInString(s string, config Config, rules ...StringReplacer) (string, error) {
	for _, rule := range rules {
		switch rule {
		case CamelCase:
			// Replace placeholders in camelCase
			withRules, err := toCamelCase(config.Metadata.PackageName)
			if err != nil {
				return "", err
			}
			s = strings.ReplaceAll(s, "{{.Metadata.PackageName}}", withRules)
		case DotsToSlashes:
			// Replace dots with slashes
			s = strings.ReplaceAll(s, ".", "/")
		default:
			return "", fmt.Errorf("unsupported string replacer: %s", rule)
		}
	}

	return s, nil
}

// toCamelCase converts a string to camelCase
func toCamelCase(s string) (string, error) {
	// Split the string by spaces
	parts := strings.Fields(s)

	// Convert each word to camelCase
	for i, part := range parts {
		if i == 0 {
			// Lowercase the first word
			parts[i] = strings.ToLower(part)
		} else {
			// Uppercase the first letter of the word
			parts[i] = strings.ToUpper(part[:1]) + part[1:]
		}
	}

	// Join the words back together
	return strings.Join(parts, ""), nil
}
