package args

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

type CommandArg struct {
	Name       string `yaml:"name"`
	Command    string `yaml:"command"`
	ExitOnFail bool   `yaml:"exit_on_fail"`
}

type CommandArgProcessor struct{}

func (p *CommandArgProcessor) Prepare(arg interface{}, cmd *cobra.Command) error {
	return nil
}

func (p *CommandArgProcessor) Process(arg interface{}, cmd *cobra.Command, templateArgs map[string]string) error {
	commandArg, ok := arg.(CommandArg)
	if !ok {
		return fmt.Errorf("invalid argument type for CommandArgProcessor")
	}
	output, err := exec.Command("sh", "-c", commandArg.Command).Output()
	if err != nil {
		if commandArg.ExitOnFail {
			return fmt.Errorf("command %s failed: %w", commandArg.Command, err)
		}
		fmt.Printf("command %s failed: %v\n", commandArg.Command, err)
	}
	templateArgs[commandArg.Name] = string(output)
	return nil
}
