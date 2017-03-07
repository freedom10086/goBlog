package handlers

import (
	"net/http"
)

// /post
func PostHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		handleGet(w, r)
	} else if (r.Method == "POST") {
		handlePost(w, r)
	}
}

func handleGet(w http.ResponseWriter, r *http.Request) {

}

func handlePost(w http.ResponseWriter, r *http.Request){

}