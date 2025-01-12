package main

import (
	"github.com/andvarfolomeev/promptflex/internal/config"
	"github.com/andvarfolomeev/promptflex/internal/snippet"

	"github.com/spf13/cobra"
)

func main() {
	config, err := config.Load()
	if err != nil {
		panic(err)
	}

	rootCmd := &cobra.Command{Use: "promptflex"}

	for _, snip := range config.Snippets {
		cmd := snippet.NewCommand(snip)
		rootCmd.AddCommand(cmd)
	}

	rootCmd.Execute()
}
