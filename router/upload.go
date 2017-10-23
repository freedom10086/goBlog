package router

import (
	"net/http"
	"time"
	"os"
	"io"
	"strings"
)

type UploadHandler struct {
	BaseHandler
}

//上传附件接口
func (*UploadHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	mod := r.FormValue("type")

	if mod != "file" && mod != "image" {
		BadParameter(w, r)
	}

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		InternalError(w, r, err)
		return
	}
	defer file.Close()

	t := time.Now()
	year := t.Year()
	month := t.Month()
	dir := config.DirUpload + "/" + string(year) + "/" + string(month) + "/"
	if err := os.MkdirAll(dir, os.ModeDir); err != nil {
		Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var fileName string
	if index := strings.Index(handler.Filename, "."); index >= 0 {
		fileName = dir + mod + "_" + string(t.Nanosecond()) + handler.Filename[:index]
	} else {
		fileName = dir + mod + "_" + string(t.Nanosecond()) + ".file"
	}

	f, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		InternalError(w, r, err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}


