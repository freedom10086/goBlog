package router

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
	"io"
	"bufio"
)

const cacheTime = time.Hour * 24
const defaultType = "application/octet-stream"

var mineTypes = map[string]string{
	".css":      "text/css; charset=utf-8",
	".json":     "application/json; charset=utf-8",
	".txt":      "text/plain; charset=utf-8",
	".gif":      "image/gif",
	".htm":      "text/html; charset=utf-8",
	".html":     "text/html; charset=utf-8",
	".jpg":      "image/jpeg",
	".jpeg":     "image/jpeg",
	".png":      "image/png",
	".svg":      "image/svg+xml",
	".ico":      "image/x-icon",
	".xml":      "text/xml; charset=utf-8",
	".js":       "application/x-javascript",
	".pdf":      "application/pdf",
	".manifest": "text/cache-manifest",
}

//默认的handler
type DefaultHandler struct {
	BaseHandler
}

func MineType(name string) string {
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

func (*DefaultHandler) DoAuth(int, *http.Request) error {
	return nil
}

func (h *DefaultHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	filename := staticDir + r.URL.Path[1:]
	if strings.HasSuffix(r.URL.Path, "/") {
		filename += "index.html"
	}
	fi, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			NotFound(w, r)
			return
		}
		const modeType = os.ModeDir | os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
		if (fi.Mode()&modeType != 0) || os.IsPermission(err) {
			Forbidden(w, r)
			return
		}
		InternalError(w, r, err)
		return
	}
	f, err := os.Open(filename)
	if err != nil {
		Forbidden(w, r)
		return
	}
	defer f.Close()
	cacheControl := fmt.Sprintf("public, max-age=%d", cacheTime/time.Second)
	mineType := MineType(path.Ext(filename))

	w.Header().Set("content-type", mineType)
	w.Header().Set("Content-Length", strconv.FormatInt(fi.Size(), 10))
	w.Header().Set("Cache-Control", cacheControl)

	w.WriteHeader(http.StatusOK)
	size := fi.Size()
	if size > 2048 {
		size = 2048
	}
	reader := bufio.NewReader(f)
	buf := make([]byte, size)
	for {
		n, err := reader.Read(buf)
		if n == 0 || err == io.EOF {
			break
		}
		if err != nil {
			InternalError(w, r, err)
			return
		}
		w.Write(buf[:n])
	}
}
