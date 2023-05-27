package utils

import (
	"bytes"
	"strconv"
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

func RenderTemplateInt(tpl *template.Template, ctx interface{}) (int, error) {
	valueStr, err := RenderTemplate(tpl, ctx)
	if err != nil {
		return 0, err
	}
	value, err := strconv.ParseInt(valueStr, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(value), nil
}

func RenderTemplateBool(tpl *template.Template, ctx interface{}) (bool, error) {
	valueStr, err := RenderTemplate(tpl, ctx)
	if err != nil {
		return false, err
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, err
	}

	return value, nil
}
