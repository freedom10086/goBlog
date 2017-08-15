package router

import (
	"net/http"
	"path"
	"strings"
	"sync"
	"log"
)

//url路径替换字符
var htmlReplacer = strings.NewReplacer(
	"&", "&amp;",
	"<", "&lt;",
	">", "&gt;",
	`"`, "&#34;",
	"'", "&#39;",
)

//路由列表
var mux map[string]Handler

//http method
const (
	MethodGet    = iota
	MethodPost
	MethodDelete
	MethodUpdate
	MethodOther
)

type MyRouter struct {
	mu             sync.RWMutex
	m              map[string]muxEntry
	defaultHandler Handler //默认handler
}

func (mux *MyRouter) DefaultHandler(h Handler) {
	mux.defaultHandler = h
}

//pattern /static->path /static/->目录
type muxEntry struct {
	h       Handler
	pattern string
}

func NewRouter() *MyRouter {
	r := new(MyRouter)

	for i, v := range mux {
		r.Register(i, v)
	}
	return r
}

//http 请求会到这儿处理
func (mux *MyRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("method:%s path:%s", r.Method, r.URL.Path)

	var m int
	switch r.Method {
	case http.MethodGet:
		m = MethodGet
	case http.MethodPost:
		m = MethodPost
	case http.MethodDelete:
		m = MethodDelete
	case http.MethodPut, http.MethodPatch:
		m = MethodUpdate
	default:
		m = MethodOther
	}

	var p string
	if p = cleanPath(r.URL.Path); p != r.URL.Path {
		log.Printf("Redirect:%s", p)
		Redirect(w, r, p, http.StatusMovedPermanently)
		return
	}

	if h := mux.handle(p); h != nil {
		err := h.DoAuth(m, r)
		if err != nil {
			Unauthorized(w, r)
			return
		}
		if m == MethodGet {
			h.DoGet(w, r)
		} else if m == MethodPost {
			h.DoPost(w, r)
		} else if m == MethodDelete {
			h.DoDelete(w, r)
		} else if m == MethodUpdate {
			h.DoUpdate(w, r)
		} else {
			NotAllowed(w, r)
			return
		}
	} else { //defaultHandler==nil 说明默认的handler未设置
		NotAllowed(w, r)
		return
	}
}

//注册路由
func (mux *MyRouter) Register(pattern string, handler Handler) {
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

func (mux *MyRouter) handle(path string) (h Handler) {
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

	if h == nil { //最终还是nil 启用默认的handler
		h = mux.defaultHandler
	}
	return
}

//isdir 表示pattern是否为目录
func pathMatch(pattern, path string) bool {
	if pattern[len(pattern)-1] != '/' || pattern == "/" {
		//如果不是目录比较是否相等
		return pattern == path
	}

	//如果是目录 则只要目录部分相同就行
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
