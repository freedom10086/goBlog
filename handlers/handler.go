package handlers

import (
	"fmt"
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

func HandleError(err error,w http.ResponseWriter, r *http.Request)  {
	w.Write([]byte(err.Error()))
}

