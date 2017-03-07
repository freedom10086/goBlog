package main

import (
	"goBlog/handlers"
	"net/http"
)

var routers map[string]http.Handler

func init() {
	readConfig()
	routers = make(map[string]http.Handler)
	routers["/"] = &handlers.StaticFileHandler{}
	routers["/cate"] = &handlers.CateHandler{}
	routers["/user"] = &handlers.UserHandler{}
	routers["/login"] = &handlers.LoginHandler{SecretKey: config.SecretKey}
}
