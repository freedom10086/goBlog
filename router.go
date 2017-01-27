package main

import (
	"goweb/controls"
	"net/http"
)

type Route struct {
	Path string
	Func func(w http.ResponseWriter, r *http.Request)
}

// 路由规则
var Routes = []Route{
	{"/", controls.HomeHandler},
	{"/user", controls.UserHandler},
	{"/post", controls.PostHandler},
}
