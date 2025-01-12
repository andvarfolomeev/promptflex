package args

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type EnvArg struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
	Variable string `yaml:"variable"`
}

type EnvArgProcessor struct{}

func (p *EnvArgProcessor) Prepare(arg interface{}, cmd *cobra.Command) error {
	return nil
}

func (p *EnvArgProcessor) Process(arg interface{}, cmd *cobra.Command, templateArgs map[string]string) error {
	envArg, ok := arg.(EnvArg)
	if !ok {
		return fmt.Errorf("invalid argument type for EnvArgProcessor")
	}
	value, ok := os.LookupEnv(envArg.Variable)
	if !ok && envArg.Required {
		return fmt.Errorf("failed getting %s environment variable", envArg.Name)
	}
	templateArgs[envArg.Name] = value
	return nil
}
