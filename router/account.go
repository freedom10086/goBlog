package router

import (
	"net/http"
	"strings"
)

//个人账号页面
type AccountHandler struct {
	BaseHandler
}

func (h *AccountHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (*AccountHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	p := r.URL.Path

	index := strings.LastIndex(p, "/")
	if index > 0 {
		switch p[index+1:] {
		case "reset_password":
			StaticTemplate(w, &TemplateData{
				Css: []string{"style.css"},
				Js:  []string{"base.js", "particles.js"},
			}, "reset_password")
		case "setting":
			StaticTemplate(w, &TemplateData{
				Css: []string{"style.css"},
			}, "setting")
		default:
			NotFound(w, r)
		}

	} else { //个人账户 todo
		NotFound(w, r)
	}
}
