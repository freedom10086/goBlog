package router

import (
	"goBlog/repository"
	"io"
	"log"
	"net/http"
	"strconv"
)

type CategoryHandler struct {
	BaseHandler
}

type CateGoryTemplateData struct {
	BasePageData
	Categories []*repository.Category
}

func (*CategoryHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (*CategoryHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	cates, err := repository.GetCates()
	if err != nil {
		InternalError(w, r, err)
		return
	}
	Template(w,
		&TemplateData{
			Title: "分类-" + config.SiteName,
			Css:   []string{"style.css"},
			Js:    []string{"base.js"},
			Data: &CateGoryTemplateData{
				BasePageData: BasePageData{
					TabIndex: 1,
				},
				Categories: cates,
			},
		},
		"page.gohtml", "category.gohtml")
}

type CategoryApiHandler struct {
	BaseHandler
}

func (*CategoryApiHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	if cates, err := repository.GetCates(); err != nil {
		InternalError(w, r, err)
	} else {
		Result(w, r, cates)
	}
}

func (*CategoryApiHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	des := r.PostFormValue("description")

	if name == "" || des == "" {
		BadParameter(w, r, "参数不足")
		return
	}

	log.Printf("name:%s des:%s", name, des)
	if i, err := repository.AddCate(name, des); err != nil {
		InternalError(w, r, err)
		return
	} else {
		log.Printf("insert cate %d ok", i)
		Result(w, r, i)
	}
}

func (*CategoryApiHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
	cid := r.PostFormValue("cid")
	if cidInt, err := strconv.Atoi(cid); err != nil {
		BadParameter(w, r, err.Error())
	} else {
		if i, err := repository.DelCate(cidInt); err != nil {
			InternalError(w, r, err)
		} else {
			log.Printf("delete cate %d ok,delete count %d", cidInt, i)
			Result(w, r, i)
		}
	}
}

func (*CategoryApiHandler) DoUpdate(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "CategoryApiHandler DoUpdate")
}
