package router

import "net/http"

type ArticleHandler struct {
	BaseHandler
}

func (h *ArticleHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (h *ArticleHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w, &TemplateData{
		Css: []string{"style.css"},
		Js:  []string{"base.js", "highlight.pack.js", "marked.min.js"}},
		"page.gohtml", "post.gohtml",
	)
}
