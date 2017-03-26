package router

import (
	"net/http"
	"path"
	"strings"
	"sync"
	"log"
)

var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

var routers map[string]MyHandler

const (
	MethodGet    = iota
	MethodPost
	MethodDelete
	MethodUpdate
	MethodOther
)

func init() {
	routers = make(map[string]MyHandler)
	routers["/static/"] = &StaticFileHandler{}
	routers["/categorys"] = &CateHandler{}
	routers["/users"] = &UserHandler{}
	routers["/auth"] = &OauthHandler{}
	routers["/register"] = &RegisterHandler{}
	routers["/chats"] = &ChatHandler{}
}

type MyRouter struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	h       MyHandler
	pattern string //pattern /static->path /static/->目录
}

//子类要实现此接口中的方法如果不实现
//由父类代替
//int HttpMethod
type MyHandler interface {
	DoAuth(int, *http.Request) error
	DoGet(http.ResponseWriter, *http.Request)
	DoPost(http.ResponseWriter, *http.Request)
	DoDelete(http.ResponseWriter, *http.Request)
	DoUpdate(http.ResponseWriter, *http.Request)
}

func NewRouter() *MyRouter {
	r := new(MyRouter)

	for i, v := range routers {
		r.Handle(i, v)
	}
	return r
}

func (mux *MyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("method:%s path:%s", r.Method, r.URL.Path)
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//代理
	if r.Method == "CONNECT" {
		NotAllowed(w, r)
		return
	}

	var p string
	if p = cleanPath(r.URL.Path); p != r.URL.Path {
		log.Printf("Redirect:%s", p)
		Redirect(w, r, p, http.StatusMovedPermanently)
		return
	}

	if h := mux.handle(p); h != nil {
		var method int
		switch r.Method {
		case "GET":
			method = MethodGet
			err := h.DoAuth(method, r)
			if err != nil {
				Unauthorized(w, r)
				return
			}
			h.DoGet(w, r)
		case "POST":
			method = MethodPost
			err := h.DoAuth(method, r)
			if err != nil {
				Unauthorized(w, r)
				return
			}
			h.DoPost(w, r)
		case "DELETE":
			method = MethodDelete
			err := h.DoAuth(method, r)
			if err != nil {
				Unauthorized(w, r)
				return
			}
			h.DoDelete(w, r)
		case "PUT":
		case "PATCH":
			method = MethodUpdate
			err := h.DoAuth(method, r)
			if err != nil {
				Unauthorized(w, r)
				return
			}
			h.DoUpdate(w, r)
		default:
			method = MethodOther
			NotAllowed(w, r)
		}
		return
	}

	NotFound(w, r)
}

//暴露自定义MyHandler接口
func (mux *MyRouter) Handle(pattern string, handler MyHandler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()
	n := len(pattern)
	if n == 0 || pattern[0] != '/' {
		panic("http: pattern shou be not null and start with / for:" + pattern)
	}
	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}
	if _, ok := mux.m[pattern]; ok {
		panic("http: multiple registrations for " + pattern)
	}
	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}
}

func (mux *MyRouter) handle(path string) (h MyHandler) {
	mux.mu.RLock()
	defer mux.mu.RUnlock()
	var maxn = 0
	for pattern, mux := range mux.m {
		//找到匹配最长的url
		//如 有 /static/ 和 /static/js/ 2个pattern
		//path为/static/js/my.js将会匹配到后一个
		if !pathMatch(pattern, path) {
			continue
		}
		n := len(pattern)
		if h == nil || n > maxn {
			maxn = n
			h = mux.h
		}
	}
	return
}

//isdir 表示pattern是否为目录
func pathMatch(pattern, path string) bool {
	if pattern[len(pattern)-1] != '/' {
		//如果不是目录比较是否相等
		return pattern == path
	}
	//是目录则前部分匹配即可
	return strings.HasPrefix(path, pattern)
}

//最少也要返回/
func cleanPath(p string) string {
	if p == "" {
		return "/"
	} else if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}
