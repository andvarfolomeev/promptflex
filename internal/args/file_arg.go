package args

import (
	"fmt"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/spf13/cobra"
)

type FileArg struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
}

type FileArgProcessor struct{}

func (p *FileArgProcessor) Prepare(arg interface{}, cmd *cobra.Command) error {
	return nil
}

func (p *FileArgProcessor) Process(arg interface{}, cmd *cobra.Command, templateArgs map[string]string) error {
	fileArg, ok := arg.(FileArg)
	if !ok {
		return fmt.Errorf("invalid argument type for FileArgProcessor")
	}

	resolvedPath, err := resolveHomePath(fileArg.Path)
	if err != nil {
		return err
	}

	file, err := os.Open(resolvedPath)
	if err != nil {
		return fmt.Errorf("Failed opening %s file for %s argument", resolvedPath, fileArg.Name)
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Failed reading %s file for %s argument", resolvedPath, fileArg.Name)
	}

	templateArgs[fileArg.Name] = string(fileData)
	return nil
}

func resolveHomePath(path string) (string, error) {
	if path[:2] != "~/" {
		return path, nil
	}
	usr, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed getting user home directory: %v", err)
	}
	return filepath.Join(usr.HomeDir, path[2:]), nil
}
