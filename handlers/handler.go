package handlers

import (
	"time"
	"net/http"
	"path"
	"os"
	"errors"
	"fmt"
	"mime"
	"strconv"
	"io/ioutil"
	"log"
	"strings"
)

type StaticServer struct {
	Dir       string
	MaxAge    time.Duration
	MIMETypes map[string]string
}

var mineTypes map[string]string = map[string]string{
	".css":  "text/css; charset=utf-8",
	".js":   "text/javascript; charset=utf-8",
	".json": "application/json; charset=utf-8",
	".txt":  "text/plain; charset=utf-8",
}

func NewStaticServer() *StaticServer {
	return &StaticServer{
		Dir:      "view",
		MaxAge:   time.Hour,
		MIMETypes:mineTypes,
	}
}

func ErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(err.Error()))
}

func (s *StaticServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fname := s.Dir + r.URL.Path

	if strings.HasSuffix(fname, "/") {
		fname += "index.html"
	}

	log.Println(fname)

	ext := path.Ext(fname)
	var mimeType string

	f, err := os.Open(fname)

	if err != nil {
		ErrorHandler(err, w, r)
		return
	}

	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		ErrorHandler(err, w, r)
		return
	}

	const modeType = os.ModeDir | os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
	if fi.Mode()&modeType != 0 {
		ErrorHandler(errors.New("not a regular file"), w, r)
		return
	}

	cacheControl := fmt.Sprintf("public, max-age=%d", s.MaxAge/time.Second)

	mimeType = s.MIMETypes[ext]

	if mimeType == "" {
		mimeType = mime.TypeByExtension(ext)
	}

	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	w.Header().Set("content-type", mimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	w.Header().Set("Cache-Control", cacheControl)

	fd, _ := ioutil.ReadAll(f)
	f.Close();
	w.WriteHeader(http.StatusOK)

	if r.Method != "HEAD" {
		w.Write(fd)
	}
}
