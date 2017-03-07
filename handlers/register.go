package handlers

import "net/http"

type RegisterHandler struct {
	BaseHandler
}

func (h *RegisterHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	HandleMethod(h, w, r)
}


