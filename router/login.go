package router

import (
	"net/http"
	"fmt"
	"time"
)

type LoginHandler struct {
	BaseHandler
}

//二维码登录
type QrLoginHandler struct {
	BaseHandler
}

func (h *LoginHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//token null /regiest ->登陆页面
//登陆页面 ->dopost -> 发邮件 -> 点击连接 -> user.doPost 插入数据库
func (h *LoginHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	StaticTemplate(w, &TemplateData{
		Css: []string{"bootstrap.css"},
		Js:  []string{"base.js", "particles.js"}, }, "login")
}

//填好用户名 邮箱假注册
//真注册链接在邮件
func (h *LoginHandler) DoPost(w http.ResponseWriter, r *http.Request) {

}

func (h *QrLoginHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//利用html5 Server-Sent推送
func (h *QrLoginHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	//event:事件类型.
	//data:消息的数据字段.
	//id:事件ID.
	//retry:一个整数值,指定了重新连接的时间(单位为毫秒),如果该字段值不是整数,则会被忽略.
	//每个字段以\n\n结尾data要换行用\r\n
	for {
		str := "data: " + fmt.Sprint(time.Now()) + "\n\n"
		str = str + "retry:" + "1000" + "\n\n"
		w.Write([]byte(str))
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
		time.Sleep(500 * time.Millisecond)
	}
}
