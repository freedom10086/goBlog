package router

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goBlog/conf"
	"goBlog/logger"
	"goBlog/repository"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var config *conf.Config

//静态文件&模板目录
var staticDir = "static/"

const baseTmpl = "page.gohtml"

func init() {
	config = conf.Conf
	if config.DirStatic != "" {
		staticDir = config.DirStatic
	}
}

//基本api返回data
type ApiData struct {
	Data    interface{} `json:"data"`
	Code    int         `json:"code"`
	Message string      `json:"message"`
}

//基本模板data返回类型基类
//返回类型要继承
type TemplateData struct {
	Data  interface{}
	Title string
	Css   []string
	Js    []string
}

//常用auth优先级 post form > get form > head Authorization > cookie["token"]
//基础用户
func BaseAuth(method int, r *http.Request) (*repository.Token, error) {
	var auth string
	if cookie, err := r.Cookie("token"); err == nil && cookie != nil {
		auth = cookie.Value
		// may need url decode = => %3D
		if strings.Contains(cookie.Value, "%") {
			if t, err := url.QueryUnescape(cookie.Value); err == nil {
				auth = t
			}
		}
		logger.D("auth cookie token %s", auth)
	}
	if a := r.Header.Get("Authorization"); len(a) > 5 && strings.HasPrefix(a, "Base ") {
		if decodeToken, err := base64.URLEncoding.DecodeString(a[5:]); err != nil {
			return nil, repository.ErrTokenInvalid
		} else {
			auth = string(decodeToken)
		}
		logger.D("auth head token %s", auth)
	}
	if formToken := r.FormValue("token"); len(formToken) > 0 {
		logger.D("auth form token %s", formToken)
		auth = formToken
	}
	if postFormToken := r.PostFormValue("token"); len(postFormToken) > 0 {
		logger.D("auth post form token %s", postFormToken)
		auth = postFormToken
	}

	if len(auth) == 0 {
		return nil, repository.ErrTokenInvalid
	}

	return repository.ValidToken(auth, config.SecretKey)
}

func Error(w http.ResponseWriter, error string, code int) {
	logger.E("error %s code %d", error, code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	w.Write([]byte(error))
	w.Write([]byte("\n"))
}

func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	Error(w, "500 Internal Server Error:"+err.Error(), http.StatusInternalServerError)
}

func Unauthorized(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, "401 Unauthorized "+msg, http.StatusUnauthorized)
}

func Forbidden(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, "403 Forbidden "+msg, http.StatusForbidden)
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	Error(w, "404 page not found", http.StatusNotFound)
}

func NotAllowed(w http.ResponseWriter, r *http.Request) {
	Error(w, "405 Method Not Allowed", http.StatusMethodNotAllowed)
}

func BadParameter(w http.ResponseWriter, r *http.Request, msg string) {
	Error(w, "400 Bad Request "+msg, http.StatusBadRequest)
}

func Result(w http.ResponseWriter, r *http.Request, data interface{}) {
	//todo 有了http状态码 还需要status? 可以不需要包装了，直接发送完整的数据
	var res *ApiData
	if data == nil {
		res = &ApiData{
			Data:    nil,
			Code:    http.StatusNoContent,
			Message: "",
		}
	} else {
		res = &ApiData{
			Data:    data,
			Code:    http.StatusOK,
			Message: "",
		}
	}

	if b, err := json.Marshal(res); err != nil {
		InternalError(w, r, err)
		return
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Write(b)
		fmt.Println(string(b))
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

//path 要转到的url
//注意这事是站内redirect
//最终setheader例子 /index.html
func Redirect(w http.ResponseWriter, r *http.Request, path string, code int) {
	oldFullUrl := r.URL.String()
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
