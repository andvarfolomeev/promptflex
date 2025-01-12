package snippet

import (
	"fmt"
	"strings"

	"github.com/andvarfolomeev/promptflex/internal/args"
	"github.com/andvarfolomeev/promptflex/internal/openai"

	"github.com/spf13/cobra"
)

func run(cmd *cobra.Command, _ []string, snippet Snippet) {
	templateArgs, err := args.Process(snippet.Args, cmd)
	if err != nil {
		panic(err)
	}

	messages := make([]openai.ReqMessage, 0)
	for _, promptTemplate := range snippet.Prompts {
		var builder strings.Builder
		err = executeTemplate(&builder, promptTemplate.Template, templateArgs)
		if err != nil {
			panic(err)
		}
		messages = append(messages, openai.ReqMessage{Content: builder.String(), Role: promptTemplate.Role})
	}

	requestBody := openai.CompletionReq{Messages: messages, Model: "gpt-4o"}
	completion, err := openai.FetchCompletions(requestBody)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%v", completion.GetText())
}

func NewCommand(snippet Snippet) *cobra.Command {
	cmd := &cobra.Command{
		Use:   snippet.Name,
		Short: snippet.Description,
		Run: func(cmd *cobra.Command, args []string) {
			run(cmd, args, snippet)
		},
	}

	err := args.Prepare(snippet.Args, cmd)

	if err != nil {
		panic(fmt.Errorf("Failed preparing variables for %s, %w", snippet.Name, err))
	}

	return cmd
}
