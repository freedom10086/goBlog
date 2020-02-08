package router

import (
	"errors"
	"goBlog/logger"
	"goBlog/repository"
	"log"
	"net/http"
	"strconv"
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
		BadParameter(w, r, "参数:page不合法")
		return
	}

	if size := r.FormValue("size"); size == "" {
		sizeInt = 30
	} else if sizeInt, err = strconv.Atoi(size); err != nil || sizeInt <= 0 {
		BadParameter(w, r, "参数:size不合法")
		return
	}

	log.Printf("page:%d size:%d", pageInt, sizeInt)
	if users, err := repository.GetUsers(pageInt, sizeInt); err != nil {
		InternalError(w, r, err)
		return
	} else {
		Result(w, r, users)
	}
}

//要验证regtoken
//注册完成，添加用户
func (h *UserHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	token := r.PostFormValue("token")
	password := r.PostFormValue("password")
	sex := r.PostFormValue("sex")
	sexInt := -1
	if sex == "0" || sex == "" {
		sexInt = 0
	} else if sex == "1" {
		sexInt = 1
	} else if sex == "2" {
		sexInt = 2
	} else {
		BadParameter(w, r, "sex非法:"+sex)
		return
	}

	if len(token) < 32 || len(password) < 6 || sexInt < 0 || sexInt > 2 {
		BadParameter(w, r, "参数错误")
		return
	}

	if t, err := repository.ValidRegToken(token, config.SecretKey); err == nil {
		log.Println("token is ok")
		id, err := repository.AddUser(t.Username, password, t.Email, sexInt)
		if err != nil {
			InternalError(w, r, err)
			return
		}
		//todo 注册成功返回token
		log.Printf("insert user %d ok", id)
		Result(w, r, id)
		return
	} else {
		logger.E("reg token is invalid %s %v", token, err)
		Unauthorized(w, r, err.Error())
	}
}

func (*UserHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
	uid := r.PostFormValue("uid")
	if uidInt, err := strconv.Atoi(uid); err != nil {
		BadParameter(w, r, "uid不合法")
		return
	} else {
		if i, err := repository.DelUser(uidInt); err != nil {
			InternalError(w, r, err)
			return
		} else {
			log.Printf("delete user %d ok,delete count %d", uidInt, i)
			Result(w, r, i)
		}
	}
}
