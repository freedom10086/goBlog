package router

import (
	"net/http"
	"path"
	"strings"
	"sync"
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
	MethodGet = iota
	MethodPost
	MethodDelete
	MethodUpdate
	MethodOther
)

func init() {
	routers = make(map[string]MyHandler)
	routers["/"] = &StaticFileHandler{}
	routers["/categorys"] = &CateHandler{}
	routers["/users"] = &UserHandler{}
	routers["/auth"] = &OauthHandler{}
	routers["/regiest"] = &RegisterHandler{}
	routers["/chats"] = &ChatHandler{}
}

type MyRouter struct {
	mu sync.RWMutex
	m  map[string]muxEntry
}

type muxEntry struct {
	h       MyHandler
	pattern string
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
	if r.RequestURI == "*" {
		if r.ProtoAtLeast(1, 1) {
			w.Header().Set("Connection", "close")
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if h, _ := mux.getHandler(r); h != nil {
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

func (mux *MyRouter) Handle(pattern string, handler MyHandler) {
	mux.mu.Lock()
	defer mux.mu.Unlock()

	n := len(pattern)
	if n == 0 || pattern[0] != '/' {
		panic("http: pattern shou be not null and start with / for:" + pattern)
	}

	// /tree/-> /tree
	if n > 1 && pattern[n - 1] == '/' {
		pattern = pattern[:n - 1]
	}

	if mux.m == nil {
		mux.m = make(map[string]muxEntry)
	}

	if _, ok := mux.m[pattern]; ok {
		panic("http: multiple registrations for " + pattern)
	}

	mux.m[pattern] = muxEntry{h: handler, pattern: pattern}
}

func (mux *MyRouter) getHandler(r *http.Request) (h MyHandler, pattern string) {
	if r.Method != "CONNECT" {
		if p := cleanPath(r.URL.Path); p != r.URL.Path {
			r.URL.Path = p
		}
	}

	mux.mu.RLock()
	defer mux.mu.RUnlock()

	//找到匹配最长的url
	var n = 0
	for k, v := range mux.m {
		if !pathMatch(k, r.URL.Path) {
			continue
		}
		if h == nil || len(k) > n {
			n = len(k)
			h = v.h
			pattern = v.pattern
		}
	}
	return
}

func pathMatch(pattern, path string) bool {
	if len(pattern) == 0 {
		// should not happen
		return false
	}
	n := len(pattern)
	if pattern[n - 1] != '/' {
		return pattern == path
	}
	return len(path) >= n && path[0:n] == pattern
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	if p[len(p) - 1] == '/' && np != "/" {
		np += "/"
	}
	return np
}

func htmlEscape(s string) string {
	return htmlReplacer.Replace(s)
}
