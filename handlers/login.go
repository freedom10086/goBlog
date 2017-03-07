package handlers

import (
	"encoding/json"
	"goBlog/code"
	"goBlog/models"
	"io"
	"log"
	"net/http"
	"time"
)

type LoginHandler struct {
	BaseHandler
	SecretKey string
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
		if len(token) <= 32 {
			HandleError(code.ERR_TOKEN_INVALID, w, r)
			return
		}
		if uid, err := models.ValidToken(token, h.SecretKey); err == nil {
			if u, err := models.GetUserById(uid); err != nil {
				HandleError(err, w, r)
				return
			} else {
				if ub, err := json.Marshal(u); ub != nil {
					w.Write(ub)
					return
				} else if err != nil {
					HandleError(err, w, r)
					return
				}

			}
		} else {
			HandleError(err, w, r)
			return
		}
	}

	//账号密码登陆
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	email := r.PostFormValue("email")

	log.Println("login")

	if (username == "" && email == "") || password == "" {
		HandleError(code.ERR_PARAMETER, w, r)
		return
	}

	log.Println("login")

	u, err := models.UserLogin(username, email, password)
	if err != nil {
		HandleError(err, w, r)
		return
	}

	if u == nil {
		HandleError(code.ERR_LOGIN, w, r)
		return
	}

	t := models.GenToken(u.Uid, h.SecretKey, time.Hour*24*30)
	log.Printf("token is :%s", t)

	ub, _ := json.Marshal(u)

	io.WriteString(w, t)
	w.Write(ub)
}
