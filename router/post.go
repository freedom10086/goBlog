package router

import "net/http"

type PostHandler struct {
	BaseHandler
}

func (h *PostHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (h *PostHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w, &TemplateData{
		Css: []string{"style.css"},
		Js:  []string{"base.js", "highlight.pack.js", "marked.min.js"}},
		"page.gohtml", "post.gohtml",
	)
}

type NewPostHandler struct {
	BaseHandler
}

type NewPostTemplateData struct {
	BasePageData
}

func (h *NewPostHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w, &TemplateData{
		Css:  []string{"style.css", "editor.css", "font-md.css"},
		Js:   []string{"base.js", "highlight.pack.js", "marked.min.js", "editor.js"},
		Data: &NewPostTemplateData{}},
		"page.gohtml", "newpost.gohtml",
	)
}
