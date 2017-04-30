package router

import (
	"net/http"
	"goBlog/model"
)

//handler接口
//子类要么实现此接口要么继承BaseHandler
type Handler interface {
	//基础auth 子类可以重写此方法实现自定义auth
	//error为nil验证成功 否则失败
	DoAuth(method int, r *http.Request) error

	//处理get请求
	DoGet(w http.ResponseWriter, r *http.Request)

	//处理post请求
	DoPost(w http.ResponseWriter, r *http.Request)

	//处理delete请求
	DoDelete(w http.ResponseWriter, r *http.Request)

	//处理update请求
	DoUpdate(w http.ResponseWriter, r *http.Request)
}

//默认的handler
type BaseHandler struct {
	Token *model.Token
}

func (h *BaseHandler) DoAuth(method int, r *http.Request) error {
	t, err := BaseAuth(method, r)
	if err != nil {
		return err
	}
	h.Token = t
	return nil
}

func (*BaseHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
func (*BaseHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
func (*BaseHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
func (*BaseHandler) DoUpdate(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
