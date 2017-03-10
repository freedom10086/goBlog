package handlers

import (
	"fmt"
	"goBlog/code"
	"goBlog/models"
	"io"
	"net/http"
	"time"
)

type RegisterHandler struct {
	BaseHandler
	SecretKey string
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleMethod(h, w, r)
}

func (h *RegisterHandler) ServeGET(w http.ResponseWriter, r *http.Request) {
	if token := r.FormValue("token"); token == "" {
		s := fmt.Sprintln("this is reg page")
		io.WriteString(w, s)
	} else if t, ok := models.ValidRegToken(token, h.SecretKey); ok {
		//完善个人信息完成注册 post /user
		t.Sex = 0
		t.Password = "6666"
		treuToken := models.EncodeToken(t, h.SecretKey)
		s := fmt.Sprintf("username is %s email is %s true reg link is :%s", t.Username, t.Email, treuToken)
		io.WriteString(w, s)
	} else {
		HandleError(code.ERR_TOKEN_INVALID, w, r)
	}
}

func (h *RegisterHandler) ServePOST(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")

	if username == "" || email == "" {
		HandleParaError(w, r)
		return
	}

	t := models.GenRegToken(username, email, h.SecretKey, time.Minute*30)
	s := fmt.Sprintf("reg link has send to your email you shou click the link in 30min %s", t)
	io.WriteString(w, s)
}
