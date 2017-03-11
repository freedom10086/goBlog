package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

type ResultData struct {
	Data    interface{}
	Code    int
	Message string
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

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if b, err := json.Marshal(res); err != nil {
		InternalError(w, r, err)
		return
	} else {
		w.Write(b)
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

func Redirect(w http.ResponseWriter, r *http.Request, urlStr string, code int) {
	if u, err := url.Parse(urlStr); err == nil {
		if u.Scheme == "" && u.Host == "" {
			oldpath := r.URL.Path
			if oldpath == "" {
				oldpath = "/"
			}
			if urlStr == "" || urlStr[0] != '/' {
				olddir, _ := path.Split(oldpath)
				urlStr = olddir + urlStr
			}
			var query string
			if i := strings.Index(urlStr, "?"); i != -1 {
				urlStr, query = urlStr[:i], urlStr[i:]
			}
			trailing := strings.HasSuffix(urlStr, "/")
			urlStr = path.Clean(urlStr)
			if trailing && !strings.HasSuffix(urlStr, "/") {
				urlStr += "/"
			}
			urlStr += query
		}
	}

	w.Header().Set("Location", urlStr)
	w.WriteHeader(code)

	// Shouldn't send the response for POST or HEAD; that leaves GET.
	if r.Method == "GET" {
		note := "<a href=\"" + htmlEscape(urlStr) + "\">" + http.StatusText(code) + "</a>.\n"
		fmt.Fprintln(w, note)
	}
}
