package router

import (
	"goBlog/models"
	"log"
	"net/http"
	"strconv"
)

type UserHandler struct {
	BaseHandler
	SecretKey string
}

func (*UserHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	var sizeInt, pageInt int;
	var err error
	//order := r.FormValue("order")
	if page := r.FormValue("page"); page == "" {
		pageInt = 1;
	} else if pageInt, err = strconv.Atoi(page); err != nil || pageInt <= 0 {
		BadParament(w, r)
		return
	}

	if size := r.FormValue("size"); size == "" {
		sizeInt = 30;
	} else if sizeInt, err = strconv.Atoi(size); err != nil || sizeInt <= 0 {
		BadParament(w, r)
		return
	}

	offset := sizeInt * (pageInt - 1)
	log.Printf("offset:%d limit:%d", offset, sizeInt)
	if users, err := models.GetUsers(true, offset, sizeInt); err != nil {
		InternalError(w, r, err)
		return
	} else {
		Result(w, r, users)
	}
}

//要验证regtoken
func (h *UserHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	if token := r.PostFormValue("token"); token == "" {
		BadParament(w, r)
		return
	} else if t, ok := models.ValidRegToken(token, h.SecretKey); ok {
		if t.Username == "" || t.Password == "" || t.Email == "" || t.Sex < 0 {
			BadParament(w, r)
			return
		}

		id, err := models.AddUser(t.Username, t.Password, t.Email, t.Sex)
		if err != nil {
			InternalError(w, r, err)
		}

		log.Printf("insert user %d ok", id)
		Result(w, r, id)
		return
	}

	Unauthorized(w, r)
}

func (*UserHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
	uid := r.PostFormValue("uid")
	if uidInt, err := strconv.Atoi(uid); err != nil {
		BadParament(w, r)
		return
	} else {
		if i, err := models.DelUser(uidInt); err != nil {
			InternalError(w, r, err)
			return
		} else {
			log.Printf("delete user %d ok,delete count %d", uidInt, i)
			Result(w, r, i)
		}
	}
}
