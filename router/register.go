package router

import (
	"fmt"
	"goBlog/model"
	"io"
	"net/http"
	"time"
)

type RegisterHandler struct {
	BaseHandler
}

func (h *RegisterHandler) DoAuth(method int, r *http.Request) error {
	if method == MethodGet || method == MethodPost {
		return nil
	}

	return h.BaseHandler.DoAuth(method, r)
}

//token null /regiest ->登陆页面
//登陆页面 ->dopost -> 发邮件 -> 点击连接 -> user.doPost 插入数据库
func (h *RegisterHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	if token := r.FormValue("token"); token == "" {
		s := fmt.Sprintln("this is reg page 用户名 邮件 ->提交 doPost")
		io.WriteString(w, s)
	} else if t, ok := model.ValidRegToken(token, config.SecretKey); ok {
		//完善个人信息完成注册 post /user
		t.Sex = 0
		t.Password = "6666"
		treuToken := model.EncodeToken(t, config.SecretKey)
		s := fmt.Sprintf("username is %s email is %s true reg link is :%s", t.Username, t.Email, treuToken)
		io.WriteString(w, s)
	} else {
		Unauthorized(w, r)
	}
}

//填好用户名 邮箱假注册
//真注册链接在邮件
func (h *RegisterHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	email := r.PostFormValue("email")

	if username == "" || email == "" {
		BadParament(w, r)
		return
	}

	//todo 验证合法

	token := model.GenRegToken(username, email, config.SecretKey, time.Minute * 30)

	go func() {
		content := fmt.Sprintf("欢迎你注册%s,请点击以下链接来验证你的邮箱,请在%d分钟内完成验证\r\n <a href=\"%s\">点击这儿</a>",
			config.SiteName,
			30,
			config.SiteAddr + config.SitePortSSL + "/regiest?token=" + token,
		)
		model.SendMail(email, "验证你的注册邮件", content)
	}()

	s := fmt.Sprintf("注册确认链接已经发送到你的邮箱%s" + token)
	io.WriteString(w, s)
}
