package router

import (
	"net/http"
)

type AboutHandler struct {
	BaseHandler
}

type AboutTemplateData struct {
	BasePageData
}

func (*AboutHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (*AboutHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w,
		&TemplateData{
			Title: "关于-" + config.SiteName,
			Css:   []string{"style.css"},
			Js:    []string{"base.js"},
			Data: &AboutTemplateData{
				BasePageData: BasePageData{TabIndex: 2},
			},
		},
		"page.gohtml", "about.gohtml")
}
