package snippet

import (
	"io"
	"text/template"
)

func executeTemplate(writter io.Writer, templateText string, templateArg map[string]string) error {
	t := template.New("t")
	t, err := t.Parse(templateText)
	if err != nil {
		return err
	}
	err = t.Execute(writter, templateArg)
	return err
}
