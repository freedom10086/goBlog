package handlers

import (
	"encoding/json"
	"goBlog/models"
	"io"
	"log"
	"net/http"
	"strconv"
	"goBlog/code"
)

type CateHandler struct {
	BaseHandler
}

func (h *CateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleMethod(h, w, r)
}

func (h *CateHandler) ServeGET(w http.ResponseWriter, r *http.Request) {
	cates, err := models.GetCates()

	if err != nil {
		HandleError(err, w, r)
		return
	}

	d, err := json.Marshal(cates)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	w.Write(d)
}

func (h *CateHandler) ServePOST(w http.ResponseWriter, r *http.Request) {
	name := r.PostFormValue("name")
	des := r.PostFormValue("description")

	if name == "" || des == "" {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	}

	log.Printf("name:%s des:%s", name, des)
	i, err := models.AddCate(name, des)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	log.Printf("insert cate %d ok", i)
	io.WriteString(w, string(i))
}

func (h *CateHandler) ServeDELETE(w http.ResponseWriter, r *http.Request) {
	cid := r.PostFormValue("cid")
	if cidInt, err := strconv.Atoi(cid); err != nil {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	} else {
		i, err := models.DelCate(cidInt)
		if err != nil {
			HandleError(err, w, r)
			return
		}

		log.Printf("delete cate %d ok,delete count %d", cidInt, i)
	}
}
