package controls

import (
	"io"
	"net/http"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello this is user!!")
}
