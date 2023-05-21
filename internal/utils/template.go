package utils

import (
	"bytes"
	"text/template"
)

func RenderTemplate(tpl *template.Template, ctx interface{}) (string, error) {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, ctx)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
