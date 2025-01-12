package args

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

type FlagFileArg struct {
	Name string `yaml:"name"`
	Flag string `yaml:"flag"`
}

type FlagFileArgProcessor struct{}

func (p *FlagFileArgProcessor) Prepare(arg interface{}, cmd *cobra.Command) error {
	flagArg, ok := arg.(FlagFileArg)
	if !ok {
		return fmt.Errorf("invalid argument type for FlagFileArgProcessor")
	}
	cmd.Flags().String(flagArg.Flag, "your value", "123")
	return nil
}

func (p *FlagFileArgProcessor) Process(arg interface{}, cmd *cobra.Command, templateArgs map[string]string) error {
	flagFileArg, ok := arg.(FlagFileArg)
	if !ok {
		return fmt.Errorf("invalid argument type for FlagFileArgProcessor")
	}
	path, err := cmd.Flags().GetString(flagFileArg.Flag)
	if err != nil {
		return fmt.Errorf("failed getting %s flag", flagFileArg.Name)
	}

	resolvedPath, err := resolveHomePath(path)
	if err != nil {
		return err
	}

	file, err := os.Open(resolvedPath)
	if err != nil {
		return fmt.Errorf("Failed opening %s file for %s argument", resolvedPath, flagFileArg.Name)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Failed reading %s file for %s argument", resolvedPath, flagFileArg.Name)
	}

	templateArgs[flagFileArg.Name] = string(fileData)
	return nil
}
