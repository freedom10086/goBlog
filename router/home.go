package router

import (
	"net/http"
)

type HomeHandler struct {
	BaseHandler
}

func (h *HomeHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//w http.ResponseWriter, data interface{}, res []string, tmpls ...string
func (*HomeHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w,
		&TemplateData{
			Title: "首页-" + config.SiteName,
			Css:   []string{"style.css"},
			Js:    []string{"base.js"},
			Data:  nil},
		"page.gohtml", "index.gohtml")
}
