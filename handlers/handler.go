package handlers

import (
	"encoding/json"
	"fmt"
	"goBlog/code"
	"io"
	"net/http"
)

type RedirectHandler struct {
	Url  string
	Port string
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := h.Url + h.Port + r.RequestURI
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}

type BaseHandler struct {
}

type MethodHandler interface {
	ServeGET(http.ResponseWriter, *http.Request)
	ServePOST(http.ResponseWriter, *http.Request)
	ServeDELETE(http.ResponseWriter, *http.Request)
	ServePUT(http.ResponseWriter, *http.Request)
	ServePATCH(http.ResponseWriter, *http.Request)
	ServeHEAD(http.ResponseWriter, *http.Request)
}

type Result struct {
	Data    interface{}
	Code    int
	Message string
}

func (h *BaseHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleMethod(h, w, r)
}

func (*BaseHandler) ServeGET(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Method:%s Url:%s", r.Method, r.URL.Path)
	io.WriteString(w, s)
}

func (*BaseHandler) ServePOST(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Method:%s Url:%s", r.Method, r.URL.Path)
	io.WriteString(w, s)
}

func (*BaseHandler) ServeDELETE(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Method:%s Url:%s", r.Method, r.URL.Path)
	io.WriteString(w, s)
}

func (*BaseHandler) ServePUT(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Method:%s Url:%s", r.Method, r.URL.Path)
	io.WriteString(w, s)
}

func (*BaseHandler) ServePATCH(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Method:%s Url:%s", r.Method, r.URL.Path)
	io.WriteString(w, s)
}

func (*BaseHandler) ServeHEAD(w http.ResponseWriter, r *http.Request) {
	s := fmt.Sprintf("Method:%s Url:%s", r.Method, r.URL.Path)
	io.WriteString(w, s)
}

func HandleMethod(target MethodHandler, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		target.ServeGET(w, r)
	case "POST":
		target.ServePOST(w, r)
	case "DELETE":
		target.ServeDELETE(w, r)
	case "PUT":
		target.ServePUT(w, r)
	case "PATCH":
		target.ServePATCH(w, r)
	case "HEAD":
		target.ServeHEAD(w, r)
	}
}

func HandleResult(data interface{},w http.ResponseWriter, r *http.Request) {
	res := &Result{Data: data, Code: code.CODE_OK, Message: ""}
	if b, err := json.Marshal(res); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}

func HandleParaError(w http.ResponseWriter, r *http.Request) {
	HandleError(code.ERR_PARAMETER, w, r)
}

func HandleError(err error, w http.ResponseWriter, r *http.Request) {
	res := &Result{Data: nil, Code: code.CODE_ERROR, Message: err.Error()}
	if b, err := json.Marshal(res); err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(b)
	}
}
