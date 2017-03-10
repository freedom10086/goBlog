package handlers

import (
	"goBlog/models"
	"io"
	"net/http"
	"time"
)

type LoginHandler struct {
	BaseHandler
	SecretKey string
}

type LoginResult struct {
	User  *models.User
	Token string
}

func (h *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleMethod(h, w, r)
}

func (*LoginHandler) ServeGET(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "this is login page to do")
}

//login 分2种
// 1.通过token登陆
// 2.通过用户名密码登陆 同时产生新的token
func (h *LoginHandler) ServePOST(w http.ResponseWriter, r *http.Request) {
	if token := r.PostFormValue("token"); token != "" {
		if uid, err := models.ValidLoginToken(token, h.SecretKey); err != nil {
			HandleError(err, w, r)
			return
		} else if u, err := models.GetUserById(uid); err != nil {
			HandleError(err, w, r)
			return
		} else {
			HandleResult(u, w, r)
			return
		}
	}

	//账号密码登陆
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")

	if (username == "" && email == "") || password == "" {
		HandleParaError(w, r)
		return
	} else if u, err := models.UserLogin(username, email, password); err != nil {
		HandleError(err, w, r)
		return
	} else {
		t := models.GenLoginToken(u.Uid, h.SecretKey, time.Hour*24*30)
		res := &LoginResult{User: u, Token: t}
		HandleResult(res, w, r)
		return
	}
}
