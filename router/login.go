package router

import (
	"net/http"
)


type LoginHandler struct {
	BaseHandler
}

func (h *LoginHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//token null /regiest ->登陆页面
//登陆页面 ->dopost -> 发邮件 -> 点击连接 -> user.doPost 插入数据库
func (h *LoginHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w, "login", nil)
}

//填好用户名 邮箱假注册
//真注册链接在邮件
func (h *LoginHandler) DoPost(w http.ResponseWriter, r *http.Request) {

}
