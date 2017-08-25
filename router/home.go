package router

import (
	"net/http"
	"fmt"
)

// /
type HomeHandler struct {
	BaseHandler
}

func (h *HomeHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//w http.ResponseWriter, data interface{}, res []string, tmpls ...string
func (*HomeHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.Path)
	Template(w,
		&TemplateData{
			Css:  []string{"style.css"},
			Js:   []string{"base.js"},
			Data: nil,},
		"page.tmpl","index.tmpl")
}
