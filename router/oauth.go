package router

import (
	"bytes"
	"encoding/base64"
	"goBlog/conf"
	"goBlog/model"
	"log"
	"net/http"
	"time"
)

type OauthHandler struct {
	BaseHandler
}

//用来产生token,需要[用户名:密码]base64编码
func (h *OauthHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	if auth := r.Header.Get("Authorization"); auth != "" {
		if decodeAuth, err := base64.URLEncoding.DecodeString(auth); err != nil {
			log.Println(err)
			BadParament(w, r)
			return
		} else if index := bytes.IndexByte(decodeAuth, ':'); index <= 0 {
			Unauthorized(w, r)
			return
		} else {
			userName := string(decodeAuth[:index])
			passWord := string(decodeAuth[index+1:])

			if t, err := model.GenToken(userName, passWord, 1,
				conf.Conf.SecretKey, time.Hour*24*30); err != nil {
				Unauthorized(w, r)
				return
			} else {
				Result(w, r, t)
				return
			}
		}
	}

	Unauthorized(w, r)
}
