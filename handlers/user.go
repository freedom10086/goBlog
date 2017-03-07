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

type UserHandler struct {
	BaseHandler
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

func (*UserHandler) ServePOST(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")
	sex := r.PostFormValue("sex")

	var sexInt int

	if sex == "0" {
		sexInt = 0
	} else if sex == "1" {
		sexInt = 1
	} else if sex == "2" {
		sexInt = 2
	} else {
		sexInt = -1
	}

	if username == "" || password == "" || email == "" || sexInt < 0 {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	}

	id, err := models.AddUser(username, password, email, sexInt)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	log.Printf("insert user %d ok", id)
	io.WriteString(w, string(id))
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
