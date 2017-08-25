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
	"log"
	"os"
	"time"
	"io/ioutil"
	"strconv"
)

var config *conf.Config
//静态文件&模板目录
var staticDir = "static/"

const baseTmpl = "page.tmpl"

func init() {
	config = conf.Conf
	if config.DirStatic != "" {
		staticDir = config.DirStatic
	}
}

//基本api返回data
type ApiData struct {
	Data    interface{}
	Code    int
	Message string
}

//基本模板data返回类型基类
//返回类型要继承
type TemplateData struct {
	Data interface{}
	Css  []string
	Js   []string
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
	res := &ApiData{
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

//todo 当作变量存在内存 http2 push相关文件css/js/图片等
//res css 文件或者js文件或者其他资源文件
func Template(w http.ResponseWriter, data *TemplateData, tmpls ...string) {
	httpPush(w, data)
	var ts []string
	for _, v := range tmpls {
		ts = append(ts, staticDir+v)
	}

	t, err := template.ParseFiles(ts...)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = t.Execute(w, data)
	if err != nil {
		Error(w, err.Error(), http.StatusInternalServerError)
	}
}


//static html
//静态的html 如登陆注册等页面
func StaticTemplate(w http.ResponseWriter, data *TemplateData, file string) {
	var filename string
	if strings.HasSuffix(file, ".html") {
		filename = staticDir + file
	} else {
		filename = staticDir + file + ".html"
	}

	fi, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			NotFound(w, nil)
			return
		}
		InternalError(w, nil, err)
		return
	}

	f, _ := os.Open(filename)
	defer f.Close()
	d, err := ioutil.ReadFile(filename)
	if err != nil {
		InternalError(w, nil, err)
	}
	httpPush(w, data)
	w.Write(d)
	cacheControl := fmt.Sprintf("public, max-age=%d", cacheTime/time.Second)
	w.Header().Set("content-type", "text/html")
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	w.Header().Set("Cache-Control", cacheControl)
	w.WriteHeader(http.StatusOK)
}

//http2 server push
func httpPush(w http.ResponseWriter, data *TemplateData) {
	//http2 push
	pusher, ok := w.(http.Pusher)
	if ok { // 支持http push
		//push css
		for _, v := range data.Css {
			if err := pusher.Push("/styles/"+v, nil); err != nil {
				log.Printf("Failed to push css: %v", err)
			}
		}

		for _, v := range data.Js {
			if err := pusher.Push("/js/"+v, nil); err != nil {
				log.Printf("Failed to push js: %v", err)
			}
		}
	}
}

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, "500 Internal Server Error:"+err.Error(), http.StatusInternalServerError)
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

func BadParameter(w http.ResponseWriter, r *http.Request) {
	Error(w, "400 Bad Request", http.StatusBadRequest)
}

//path 要转到的url
//注意这事是站内redirect
//最终setheader例子 /index.html
func Redirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	oldFullUrl := r.URL.String();
	if i := strings.Index(oldFullUrl, "?"); i != -1 {
		path += oldFullUrl[i:] //加上query参数
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
