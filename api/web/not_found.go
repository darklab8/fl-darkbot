package web

import (
	"darkbot/settings/utils"
	_ "embed"
	"net/http"
	"text/template"
)

//go:embed 404.md
var NotFoundPage string
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

type NotFoundException struct{}

func (n NotFoundException) Error() string { return "Not found" }

func (e *endpoint) Check404(w http.ResponseWriter, r *http.Request) error {
	if r.URL.Path != string(e.route) {
		NotFoundHandler(w, r, e.server.router)
		return NotFoundException{}
	}
	return nil
}
