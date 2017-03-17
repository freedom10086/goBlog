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

type CompeteRegData struct {
	PostUrl  string
	Token    string
	Username string
	Email    string
}

//token null /regiest ->登陆页面
//登陆页面 ->dopost -> 发邮件 -> 点击连接 -> user.doPost 插入数据库
func (h *RegisterHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	mode := r.FormValue("mod")
	switch mode {
	case "competeRegiest":
		token := r.FormValue("token")
		if t, ok := model.ValidRegToken(token, config.SecretKey); ok {
			//返回完善信息页面,完善成功后
			//post /users 插入数据库完成注册
			Template(w, "true_register", &CompeteRegData{
				PostUrl:"/users",
				Token:token,
				Email:t.Email,
				Username:t.Username})
		} else {
			Unauthorized(w, r)
		}
		return
	case "checkUsername":
		if u := r.FormValue("username"); !model.CheckUsername(u) {
			Error(w, u + "用户名不可用", 400)
		} else {
			io.WriteString(w, "ok")
		}
		return
	case "checkEmail":
		if e := r.FormValue("email"); !model.CheckEmail(e) {
			Error(w, e + "邮箱不可用", 400)
		} else {
			io.WriteString(w, "ok")
		}
		return
	default:
		Template(w, "register", nil)
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

	if !model.CheckUsername(username) {
		Error(w, username + "用户名不可用", 400)
		return
	}

	if !model.CheckEmail(email) {
		Error(w, email + "邮件不可用", 400)
		return
	}

	token := model.GenRegToken(username, email, config.SecretKey, time.Minute * 30)
	go func() {
		content := fmt.Sprintf("欢迎你注册%s,请点击以下链接来验证你的邮箱,请在%d分钟内完成验证\r\n <a href=\"%s\">点击这儿</a>",
			config.SiteName,
			30,
			"https://" + config.SiteIpAddr + config.SitePortSSL + "/regiest?mod=competeRegiest&token=" + token,
		)
		model.SendMail(email, "验证你的注册邮件", content)
	}()
	s := fmt.Sprintf("注册确认链接已经发送到你的邮箱:%s,请在30分钟内完成验证", email)
	io.WriteString(w, s)
}
