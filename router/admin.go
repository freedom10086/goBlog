package router

import (
	"net/http"
	"strings"
)

type AdminHandler struct {
	BaseHandler
}

func (h *AdminHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (*AdminHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path
	index := strings.LastIndex(p, "/")
	var filename string
	var data *TemplateData
	if index > 0 {
		switch p[index+1:] {
		case "", "home":
			filename = "admin_home"
			data = &TemplateData{
				Css: []string{"style.css"},
				Js:  []string{"base.js", "Chart.min.js"},
			}
		case "category", "comment", "post", "user":
			filename = "admin_" + p[index+1:]
			data = &TemplateData{
				Css: []string{"style.css"},
			}
		default:
			NotFound(w, r)
			return
		}

		StaticTemplate(w, data, filename)
		return
	}

	NotFound(w, r)
}
