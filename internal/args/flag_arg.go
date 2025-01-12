package args

import (
	"fmt"

	"github.com/spf13/cobra"
)

type FlagArg struct {
	Name string `yaml:"name"`
}

type FlagArgProcessor struct{}

func (p *FlagArgProcessor) Prepare(arg interface{}, cmd *cobra.Command) error {
	flagArg, ok := arg.(FlagArg)
	if !ok {
		return fmt.Errorf("invalid argument type for FlagArgProcessor")
	}
	cmd.Flags().String(flagArg.Name, "your value", "123")
	return nil
}

func (p *FlagArgProcessor) Process(arg interface{}, cmd *cobra.Command, templateArgs map[string]string) error {
	flagArg, ok := arg.(FlagArg)
	if !ok {
		return fmt.Errorf("invalid argument type for FlagArgProcessor")
	}
	value, err := cmd.Flags().GetString(flagArg.Name)
	if err != nil {
		return fmt.Errorf("failed getting %s flag", flagArg.Name)
	}
	templateArgs[flagArg.Name] = value
	return nil
}
