package utils

import (
	"bytes"
	"darkbot/utils/logger"
	"text/template"
)

func TmpRender(templateRef *template.Template, data interface{}) string {
	header := bytes.Buffer{}
	err := templateRef.Execute(&header, data)
	logger.CheckPanic(err)
	return header.String()
}

func TmpInit(content string) *template.Template {
	var err error
	templateRef, err := template.New("test").Parse(content)
	logger.CheckPanic(err)
	return templateRef
}
