package router

import (
	"goBlog/model"
	"log"
	"net/http"
	"strconv"
	"errors"
)

type UserHandler struct {
	BaseHandler
}

func (h *UserHandler) DoAuth(method int, r *http.Request) error {
	if method == MethodPost {
		//由注册页面来的真正注册请求
		//验证注册token
		if token := r.PostFormValue("token"); token == "" {
			return errors.New("reg token needed!")
		}

		return nil
	}

	return h.BaseHandler.DoAuth(method, r)
}

func (*UserHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	var sizeInt, pageInt int
	var err error
	//order := r.FormValue("order")
	if page := r.FormValue("page"); page == "" {
		pageInt = 1
	} else if pageInt, err = strconv.Atoi(page); err != nil || pageInt <= 0 {
		BadParameter(w, r)
		return
	}

	if size := r.FormValue("size"); size == "" {
		sizeInt = 30
	} else if sizeInt, err = strconv.Atoi(size); err != nil || sizeInt <= 0 {
		BadParameter(w, r)
		return
	}

	log.Printf("page:%d size:%d", pageInt, sizeInt)
	if users, err := model.GetUsers(pageInt, sizeInt); err != nil {
		InternalError(w, r, err)
		return
	} else {
		Result(w, r, users)
	}
}

//要验证regtoken
func (h *UserHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("token")
	password := r.PostFormValue("password")
	sex := r.PostFormValue("sex")
	sexInt := -1
	if sex == "0" {
		sexInt = 0
	} else if sex == "1" {
		sexInt = 1
	} else if sex == "2" {
		sexInt = 2
	}
	if len(token) < 32 || len(password) < 6 || sexInt < 0 || sexInt > 2 {
		BadParameter(w, r)
		return
	}
	if t, ok := model.ValidRegToken(token, config.SecretKey); ok {
		log.Println("token is ok")
		id, err := model.AddUser(t.Username, password, t.Email, sexInt)
		if err != nil {
			InternalError(w, r, err)
			return
		}
		//todo 注册成功返回token
		log.Printf("insert user %d ok", id)
		Result(w, r, id)
		return
	}

	log.Println("[[[[[]]]]]]")

	Unauthorized(w, r)
}

func (*UserHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
	uid := r.PostFormValue("uid")
	if uidInt, err := strconv.Atoi(uid); err != nil {
		BadParameter(w, r)
		return
	} else {
		if i, err := model.DelUser(uidInt); err != nil {
			InternalError(w, r, err)
			return
		} else {
			log.Printf("delete user %d ok,delete count %d", uidInt, i)
			Result(w, r, i)
		}
	}
}
