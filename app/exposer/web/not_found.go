package web

import (
	_ "embed"
	"net/http"
	"text/template"

	"github.com/darklab8/darklab_goutils/goutils/utils"
	"github.com/darklab8/darklab_goutils/goutils/utils/utils_types"
)

//go:embed 404.md
var NotFoundPage utils_types.TemplateExpression
var NotFoundTemplate *template.Template

func init() {
	NotFoundTemplate = utils.TmpInit(NotFoundPage)
}

type TemplateNotFoundPageVars struct {
	Routes map[route]*endpoint
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request, Routes map[route]*endpoint) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(utils.TmpRender(NotFoundTemplate, TemplateNotFoundPageVars{
		Routes: Routes,
	})))
}

type ErrorNotFound struct{}

func (n ErrorNotFound) Error() string { return "Not found" }

func (e *endpoint) Check404(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != string(e.route) {
		NotFoundHandler(w, r, e.server.router)
		return ErrorNotFound{}
	}
	return nil
}
