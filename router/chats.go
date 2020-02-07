package router

import (
	"goBlog/repository"
	"log"
	"net/http"
)

type ChatHandler struct {
	BaseHandler
}

func (*ChatHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	cs, err := repository.GetChats(1, 3, 1, 20)
	if err != nil {
		log.Print(err)
	}
	Result(w, r, cs)
}
