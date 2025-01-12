package args

import (
	"fmt"

	"github.com/spf13/cobra"
)

type Args struct {
	Envs      []EnvArg      `yaml:"envs"`
	Flags     []FlagArg     `yaml:"flags"`
	Files     []FileArg     `yaml:"files"`
	FlagFiles []FlagFileArg `yaml:"flag_files"`
	Commands  []CommandArg  `yaml:"commands"`
}

type TypedArg struct {
	Arg  interface{}
	Type string
}

func (a Args) All() []TypedArg {
	var allArgs []TypedArg

	for _, env := range a.Envs {
		allArgs = append(allArgs, TypedArg{Arg: env, Type: "env"})
	}
	for _, flag := range a.Flags {
		allArgs = append(allArgs, TypedArg{Arg: flag, Type: "flag"})
	}
	for _, file := range a.Files {
		allArgs = append(allArgs, TypedArg{Arg: file, Type: "file"})
	}
	for _, flagFile := range a.FlagFiles {
		allArgs = append(allArgs, TypedArg{Arg: flagFile, Type: "flagFile"})
	}
	for _, command := range a.Commands {
		allArgs = append(allArgs, TypedArg{Arg: command, Type: "command"})
	}

	return allArgs
}

type ArgProcessor interface {
	Process(arg interface{}, cmd *cobra.Command, templateArgs map[string]string) error
}

type ArgPreparer interface {
	Prepare(arg interface{}, cmd *cobra.Command) error
}

func Process(args Args, cmd *cobra.Command) (map[string]string, error) {
	templateArgs := make(map[string]string)

	processors := map[string]ArgProcessor{
		"env":      &EnvArgProcessor{},
		"flag":     &FlagArgProcessor{},
		"file":     &FileArgProcessor{},
		"flagFile": &FlagFileArgProcessor{},
		"command":  &CommandArgProcessor{},
	}

	for _, wrappedArg := range args.All() {
		processor, ok := processors[wrappedArg.Type]
		if !ok {
			return nil, fmt.Errorf("unknown argument type: %s", wrappedArg.Type)
		}
		if err := processor.Process(wrappedArg.Arg, cmd, templateArgs); err != nil {
			return nil, err
		}
	}

	return templateArgs, nil
}

func Prepare(args Args, cmd *cobra.Command) error {
	processors := map[string]ArgPreparer{
		"env":      &EnvArgProcessor{},
		"flag":     &FlagArgProcessor{},
		"file":     &FileArgProcessor{},
		"flagFile": &FlagFileArgProcessor{},
		"command":  &CommandArgProcessor{},
	}

	for _, wrappedArg := range args.All() {
		processor, ok := processors[wrappedArg.Type]
		if !ok {
			return fmt.Errorf("unknown argument type: %s", wrappedArg.Type)
		}
		if err := processor.Prepare(wrappedArg.Arg, cmd); err != nil {
			return err
		}
	}

	return nil
}
