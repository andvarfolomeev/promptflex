package snippet

import "github.com/andvarfolomeev/promptflex/internal/args"

type Snippet struct {
	Name        string          `yaml:"name"`
	Description string          `yaml:"description"`
	Args        args.Args       `yaml:"args"`
	Prompts     []SnippetPrompt `yaml:"prompts"`
}

type SnippetPrompt struct {
	Role     string `yaml:"role"`
	Template string `yaml:"template"`
}
