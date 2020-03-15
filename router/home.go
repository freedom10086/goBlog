package router

import (
	"goBlog/repository"
	"net/http"
)

type HomeHandler struct {
	BaseHandler
}

type HomeTemplateData struct {
	BasePageData
	Posts []*repository.Post
}

func (h *HomeHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//w http.ResponseWriter, data interface{}, res []string, tmpls ...string
func (*HomeHandler) DoGet(w http.ResponseWriter, r *http.Request) {

	posts, err := repository.GetPostListOrderLastReply(1, 30)
	if err != nil {
		InternalError(w, r, err)
		return
	}

	Template(w,
		&TemplateData{
			Title: "首页-" + config.SiteName,
			Css:   []string{"style.css"},
			Js:    []string{"base.js"},
			Data: &HomeTemplateData{
				BasePageData: BasePageData{TabIndex: 0},
				Posts:        posts,
			},
		},
		"page.gohtml", "index.gohtml")
}
