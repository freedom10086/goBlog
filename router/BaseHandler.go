package router

import (
	"net/http"
)

type BaseHandler struct {

}

func (*BaseHandler)DoGet(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w,r)
}

func (*BaseHandler)DoPost(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w,r)
}

func (*BaseHandler)DoDelete(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w,r)
}

func (*BaseHandler)DoUpdate(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w,r)
}

func (*BaseHandler)DoOther(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w,r)
}
