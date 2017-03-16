package router

import (
	"net/http"
	"log"
)

type RedirectHandler struct {
	Url  string
	Port string
}

func (h *RedirectHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := h.Url + h.Port + r.RequestURI
	log.Printf("redirect from %s to %s", r.RequestURI, url)
	http.Redirect(w, r, url, http.StatusMovedPermanently)
}