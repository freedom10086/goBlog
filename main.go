package main

import (
	"net/http"
	"goBlog/handlers"
	"log"
)

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := "https://127.0.0.1" + config.SitePortSSL + req.RequestURI
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

func init() {
	log.Print("==main init==")
}

func main() {
	go func() {
		log.Println("listen on " + config.SitePort + ". Go to http://127.0.0.1:8080/")
		http.ListenAndServe(config.SitePort, &myHandler{})
	}()

	http.Handle("/", handlers.NewStaticServer())
	log.Println("listen on " + config.SitePortSSL + ". Go to https://127.0.0.1:10443/")
	err := http.ListenAndServeTLS(config.SitePortSSL, "cert.pem", "key.pem", nil)
	log.Fatal(err)
}