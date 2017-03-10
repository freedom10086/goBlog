package handlers

import (
	"encoding/json"
	"goBlog/code"
	"goBlog/models"
	"io"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	BaseHandler
	SecretKey string
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleMethod(h, w, r)
}

func (*UserHandler) ServeGET(w http.ResponseWriter, r *http.Request) {
	//order := r.FormValue("order")
	page := r.FormValue("page")
	size := r.FormValue("size")

	if page == "" {
		page = "1"
	}

	if size == "" {
		size = "30"
	}

	var sizeInt, pageInt int
	var err error

	if pageInt, err = strconv.Atoi(page); err != nil || pageInt <= 0 {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	}

	if sizeInt, err = strconv.Atoi(size); err != nil || sizeInt <= 0 {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	}

	offset := sizeInt * (pageInt - 1)

	log.Printf("offset:%d limit:%d", offset, sizeInt)

	users, err := models.GetUsers(true, offset, sizeInt)

	if err != nil {
		HandleError(err, w, r)
		return
	}

	d, err := json.Marshal(users)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	w.Write(d)
}

//要验证regtoken
func (h *UserHandler) ServePOST(w http.ResponseWriter, r *http.Request) {
	if token := r.PostFormValue("token"); token == "" {
		HandleParaError(w, r)
		return
	} else if t, ok := models.ValidRegToken(token, h.SecretKey); ok {
		if t.Username == "" || t.Password == "" || t.Email == "" || t.Sex < 0 {
			HandleParaError(w, r)
			return
		}

		id, err := models.AddUser(t.Username, t.Password, t.Email, t.Sex)
		if err != nil {
			HandleError(err, w, r)
			return
		}

		log.Printf("insert user %d ok", id)
		io.WriteString(w, string(id))
	}

	HandleError(code.ERR_TOKEN_INVALID, w, r)
}

func (*UserHandler) ServeDELETE(w http.ResponseWriter, r *http.Request) {
	uid := r.PostFormValue("uid")
	if uidInt, err := strconv.Atoi(uid); err != nil {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	} else {
		i, err := models.DelUser(uidInt)
		if err != nil {
			HandleError(err, w, r)
			return
		}

		log.Printf("delete user %d ok,delete count %d", uidInt, i)
	}
}
