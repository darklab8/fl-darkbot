package utils

import (
	"bytes"
	"darkbot/app/settings/darkbot_logus"
	"strings"
	"text/template"
)

func TmpRender(templateRef *template.Template, data interface{}) string {
	header := bytes.Buffer{}
	err := templateRef.Execute(&header, data)
	darkbot_logus.Log.CheckFatal(err, "failed to render template")
	return header.String()
}

func TmpInit(content string) *template.Template {
	funcs := map[string]any{
		"contains":  strings.Contains,
		"hasPrefix": strings.HasPrefix,
		"hasSuffix": strings.HasSuffix}

	var err error
	templateRef, err := template.New("test").Funcs(funcs).Parse(content)
	darkbot_logus.Log.CheckFatal(err, "failed to init template")
	return templateRef
}
