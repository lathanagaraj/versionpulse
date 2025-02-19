package vp

import (
	"os"

	"gopkg.in/yaml.v2"
)

type Tool struct {
	ID   string `yaml:"ID"`
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}
type Tools struct {
	Tools []Tool `yaml:"tools"`
}

func Load() (*Tools, error) {
	// Load tools.yaml file
	data, err := os.ReadFile("tools.yaml")
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML data into a slice of Tool structs
	var toolsMap Tools
	if err := yaml.Unmarshal(data, &toolsMap); err != nil {
		return nil, err
	}

	return &toolsMap, nil
}
