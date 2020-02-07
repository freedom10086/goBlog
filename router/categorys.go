package router

import (
	"goBlog/repository"
	"io"
	"log"
	"net/http"
	"strconv"
)

type CateHandler struct {
	BaseHandler
}

func (*CateHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	if cates, err := repository.GetCates(); err != nil {
		InternalError(w, r, err)
	} else {
		Result(w, r, cates)
	}
}

func (*CateHandler) DoPost(w http.ResponseWriter, r *http.Request) {
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

func (*CateHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
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

func (*CateHandler) DoUpdate(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "CateHandler DoUpdate")
}
