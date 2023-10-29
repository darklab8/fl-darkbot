package utils

import (
	"bytes"
	"darkbot/settings/utils/logger"
	"strings"
	"text/template"
)

func TmpRender(templateRef *template.Template, data interface{}) string {
	header := bytes.Buffer{}
	err := templateRef.Execute(&header, data)
	logger.CheckPanic(err)
	return header.String()
}

func TmpInit(content string) *template.Template {
	funcs := map[string]any{
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix}

	var err error
	templateRef, err := template.New("test").Funcs(funcs).Parse(content)
	logger.CheckPanic(err)
	return templateRef
}
