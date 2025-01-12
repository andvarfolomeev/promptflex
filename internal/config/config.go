package config

import (
	"fmt"
	"io"
	"os"

	"github.com/andvarfolomeev/promptflex/internal/snippet"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Snippets []snippet.Snippet `yaml:"snippets"`
}

func Load() (Config, error) {
	path := os.ExpandEnv("$HOME/.promptflex.yml")
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf("failed opening config file: %w", err)
	}
	defer file.Close()
	yamlData, err := io.ReadAll(file)
	if err != nil {
		return Config{}, fmt.Errorf("failed reading config file: %w", err)
	}
	return Parse(yamlData)
}

func Parse(yamlData []byte) (Config, error) {
	var config Config
	if err := yaml.Unmarshal(yamlData, &config); err != nil {
		return config, err
	}
	return config, nil
}
