package router

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goBlog/conf"
	"goBlog/model"
	"net/http"
	"strings"
	"html/template"
)

var config *conf.Config

const templateDir = "template/"

type BaseHandler struct {
	Token *model.Token
}

func init() {
	config = conf.Conf
}

//基础auth 子类可以重写此方法实现自定义auth
func (h *BaseHandler) DoAuth(method int, r *http.Request) error {
	t, err := BaseAuth(method, r)
	if err != nil {
		return err
	}
	h.Token = t
	return nil
}

func (*BaseHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
func (*BaseHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
func (*BaseHandler) DoDelete(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}
func (*BaseHandler) DoUpdate(w http.ResponseWriter, r *http.Request) {
	NotAllowed(w, r)
}

type ResultData struct {
	Data    interface{}
	Code    int
	Message string
}

//常用auth
//基础用户
func BaseAuth(method int, r *http.Request) (*model.Token, error) {
	var auth string
	if auth := r.Header.Get("Authorization"); auth == "" {
		return nil, model.ErrTokenInvalid
	}
	if decodeToken, err := base64.URLEncoding.DecodeString(auth); err != nil {
		return nil, model.ErrTokenInvalid
	} else {
		return model.ValidToken(string(decodeToken), config.SecretKey)
	}
}

func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	fmt.Fprintln(w, error)
}

func Result(w http.ResponseWriter, r *http.Request, data interface{}) {
	res := &ResultData{
		Data:    data,
		Code:    http.StatusOK,
		Message: "",
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if b, err := json.Marshal(res); err != nil {
		InternalError(w, r, err)
		return
	} else {
		w.Write(b)
	}
}

//todo 当作变量存在内存
func Template(w http.ResponseWriter, tmpl string, data interface{}) {
	filepath := templateDir + tmpl + ".html"
	t, err := template.ParseFiles(filepath)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, "500 Internal Server Error:" + err.Error(), http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter, r *http.Request) {
	Error(w, "401 Unauthorized", http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, r *http.Request) {
	Error(w, "403 Forbidden", http.StatusForbidden)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, "404 page not found", http.StatusNotFound)
}

func NotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func BadParament(w http.ResponseWriter, r *http.Request) {
	Error(w, "400 Bad Request", http.StatusBadRequest)
}


//path 要转到的url
//注意这事是站内redirect
//最终setheader例子 /index.html
func Redirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	oldFullUrl := r.URL.String();
	if i := strings.Index(oldFullUrl, "?"); i != -1 {
		path += oldFullUrl[i:]//加上query参数
	}

	w.Header().Set("Location", path)
	w.WriteHeader(code)

	if !strings.HasSuffix(path, "http") {
		path = r.Host + path
	}

	// Shouldn't send the response for POST or HEAD; that leaves GET.
	if r.Method == "GET" {
		note := "<a href=\"" + htmlEscape(path) + "\">" + http.StatusText(code) + "</a>.\n"
		fmt.Fprintln(w, note)
	}
}
