package router

import "net/http"

type RedirectHandler struct {
	Url  string
	Port string
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := h.Url + h.Port + r.RequestURI
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}