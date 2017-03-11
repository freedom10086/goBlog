package router

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

const viewPath = "view"
const cacheTime = time.Hour * 24
const defaultType = "application/octet-stream"

var mineTypes = map[string]string{
	".css":  "text/css; charset=utf-8",
	".json": "application/json; charset=utf-8",
	".txt":  "text/plain; charset=utf-8",
	".gif":  "image/gif",
	".htm":  "text/html; charset=utf-8",
	".html": "text/html; charset=utf-8",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".svg":  "image/svg+xml",
	".ico":  "image/x-icon",
	".xml":  "text/xml; charset=utf-8",
	".js":   "application/x-javascript",
	".pdf":  "application/pdf",
}

type StaticFileHandler struct {
	BaseHandler
}

func GetMineType(name string) string {
	position := strings.IndexByte(name, '.')
	if position >= 0 {
		ext := name[position:]
		if v := mineTypes[ext]; v != "" {
			return v
		}

		var buf [10]byte
		lower := buf[:0]
		for i := 0; i < len(ext); i++ {
			c := ext[i]
			if 'A' <= c && c <= 'Z' {
				lower = append(lower, c+('a'-'A'))
			} else {
				lower = append(lower, c)
			}
		}

		if v := mineTypes[string(lower)]; v != "" {
			return v
		}

		//调用系统mineType
		if mineType := mime.TypeByExtension(string(lower)); mineType != "" {
			return mineType
		}

	}

	return defaultType

}

func (*StaticFileHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	fname := viewPath + r.URL.Path

	if strings.HasSuffix(fname, "/") {
		fname += "index.html"
	}

	ext := path.Ext(fname)

	f, err := os.Open(fname)
	if err != nil {
		NotFound(w, r)
		return
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		InternalError(w, r, err)
		return
	}

	const modeType = os.ModeDir | os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
	if fi.Mode()&modeType != 0 {
		InternalError(w, r, err)
		return
	}

	cacheControl := fmt.Sprintf("public, max-age=%d", cacheTime/time.Second)
	mineType := GetMineType(ext)

	log.Printf("%s %s %s", r.Method, fname, mineType)

	w.Header().Set("content-type", mineType)
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	w.Header().Set("Cache-Control", cacheControl)

	fd, _ := ioutil.ReadAll(f)
	f.Close()
	w.WriteHeader(http.StatusOK)

	if r.Method != "HEAD" {
		w.Write(fd)
	}
}
