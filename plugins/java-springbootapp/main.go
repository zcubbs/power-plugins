package main

import (
	_ "embed"
	"fmt"
	"github.com/hashicorp/go-plugin"
	"github.com/zcubbs/blueprint"
	"os"
	"path/filepath"
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

	// iterate over options and set config values
	// if no value is provided, use the default value
	for _, option := range spec.Options {
		switch option.ID {
		case "java_version":
			config.JavaVersion = values[option.ID]
			if config.JavaVersion == "" {
				config.JavaVersion = option.Default
			}
		case "packaging":
			config.Packaging = values[option.ID]
			if config.Packaging == "" {
				config.Packaging = option.Default
			}
		case "springBootVersion":
			config.SpringBootVersion = values[option.ID]
			if config.SpringBootVersion == "" {
				config.SpringBootVersion = option.Default
			}
		case "dependencyManagement":
			config.DependencyManagement = values[option.ID]
			if config.DependencyManagement == "" {
				config.DependencyManagement = option.Default
			}
		case "metadataPackageName":
			config.Metadata.PackageName = values[option.ID]
			if config.Metadata.PackageName == "" {
				config.Metadata.PackageName = option.Default
			}
		case "metadataGroupId":
			config.Metadata.GroupId = values[option.ID]
			if config.Metadata.GroupId == "" {
				config.Metadata.GroupId = option.Default
			}
		case "metadataArtifactId":
			config.Metadata.ArtifactId = values[option.ID]
			if config.Metadata.ArtifactId == "" {
				config.Metadata.ArtifactId = option.Default
			}
		case "metadataName":
			config.Metadata.Name = values[option.ID]
			if config.Metadata.Name == "" {
				config.Metadata.Name = option.Default
			}
		case "metadataDescription":
			config.Metadata.Description = values[option.ID]
			if config.Metadata.Description == "" {
				config.Metadata.Description = option.Default
			}
		}
	}

	return config, nil
}

// createProjectStructure sets up the project directory and base files
func createProjectStructure(outputPath string, config Config) (string, error) {
	projectPath := filepath.Join(outputPath, "java-springbootapp")
	if err := os.MkdirAll(projectPath, 0755); err != nil {
		return "", err
	}

	// Create subdirectories
	dirs := []string{"src", "src/main", "src/main/java", "src/main/resources", "src/test", "src/test/java"}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectPath, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return "", err
		}
	}

	return projectPath, nil
}

// go:embed templates/pom.xml
var pomXmlTemplate string

// generateProjectFiles generates project files based on the parsed config
func generateProjectFiles(projectPath string, config Config) error {
	// Define the file paths and corresponding templates
	files := map[string]string{
		"pom.xml": pomXmlTemplate,
	}

	// Process each template and create files
	for filePath, tmpl := range files {
		fullPath := filepath.Join(projectPath, filePath)
		if err := processTemplate(fullPath, tmpl, config); err != nil {
			return err
		}
	}

	return nil
}

func processTemplate(filePath, tmpl string, config Config) error {
	t, err := template.New("template").Parse(tmpl)
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
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
