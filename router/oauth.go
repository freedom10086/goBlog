package router

import (
	"bytes"
	"encoding/base64"
	"goBlog/model"
	"log"
	"net/http"
	"strings"
	"time"
)

type OauthHandler struct {
	BaseHandler
}

func (*OauthHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//用来产生token,需要[用户名:密码]base64编码
func (h *OauthHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	if auth := r.Header.Get("Authorization"); auth != "" && strings.Index(auth, "Basic ") == 0 {
		var decodeAuth []byte
		var err error
		if decodeAuth, err = base64.URLEncoding.DecodeString(auth[6:]); err != nil {
			log.Println(err)
			BadParameter(w, r, err.Error())
			return
		}

		if index := bytes.IndexByte(decodeAuth, ':'); index <= 0 {
			log.Println(decodeAuth)
			Unauthorized(w, r, "")
			return
		} else {
			userName := string(decodeAuth[:index])
			passWord := string(decodeAuth[index+1:])
			log.Printf("username is %s password is %s", userName, passWord)

			if u, err := model.UserLogin(userName, passWord); err != nil {
				Unauthorized(w, r, err.Error())
				return
			} else {
				if t, err := model.GenToken(u, 1, config.SecretKey, time.Hour*24*7); err != nil {
					Unauthorized(w, r, err.Error())
					return
				} else {
					Result(w, r, t)
					return
				}
			}
		}
	}

	Unauthorized(w, r, "")
}
