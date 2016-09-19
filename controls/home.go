package controls

import (
	"io"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello this is home!!")
}
