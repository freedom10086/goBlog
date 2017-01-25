package main

type StaticServer struct {
	Dir string
	MaxAge time.Duration
	MIMETypes map[string]string
}

var dir = "view"
var mineTypes map[string]string = {
	".css": "text/css; charset=utf-8",
	".js":  "text/javascript; charset=utf-8",
	".json": "application/json; charset=utf-8",
	".txt"： "text/plain; charset=utf-8",
}

func NewStaticServer() *StaticServer {
	return &StaticServer{
		Dir:    dir,
		MaxAge: time.Hour,
		MIMETypes: mineTypes,
}


func ErrorHandler(err error, w http.ResponseWriter, r *http.Request) {
	
}


func ResourceHandler(w http.ResponseWriter, r *http.Request) {
	ext := path.Ext(fname)
	var mimeType string

	f, err := os.Open(fname)

	if err != nil {
		return ErrorHandler(err,w,r)
	}
	
	fi, err := f.Stat()
	if err != nil {
		f.Close()
		return ErrorHandler(err,w,r)
	}

	const modeType = os.ModeDir | os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice
	if fi.Mode()&modeType != 0 {
		f.Close()
		return ErrorHandler(errors.New("not a regular file"),w,r)
	}


	cacheControl := fmt.Sprintf("public, max-age=%d", maxAge/time.Second)


	mimeType = MIMETypes(ext)

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