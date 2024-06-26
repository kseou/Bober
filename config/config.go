package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents the structure of the YAML configuration file.
type Config struct {
	Project struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"project"`
	Cpp struct {
		Compiler string `yaml:"compiler"`
		Standard string `yaml:"standard"`
		Flags    string `yaml:"flags"`
	} `yaml:"cpp"`
	Sources []string `yaml:"sources"`
	Output  struct {
		Executable string `yaml:"executable"`
	} `yaml:"output"`
	Libraries []struct {
		Name   string `yaml:"name,omitempty"`
		Config string `yaml:"config,omitempty"`
	} `yaml:"libraries,omitempty"`
}

// Load reads and parses the YAML configuration file specified by filename.
func Load(filename string) (*Config, error) {
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
