package main

import (
	"goBlog/handlers"
	"goBlog/models"
	"log"
	"net/http"
)

func init() {
	models.InitDB(config.DbName, config.DbUsername, config.DbPassword)
	log.Printf("==%s started==", config.SiteName)
}

func main() {
	go func() {
		log.Printf("http listen on %s%s", config.SiteAddr, config.SitePort)
		http.ListenAndServe(config.SitePort, &handlers.RedirectHandler{
			Url:  config.SiteAddr,
			Port: config.SitePortSSL,
		})
	}()

	mux := http.NewServeMux()
	for k, v := range routers {
		mux.Handle(k, v)
	}

	log.Printf("https listen on %s%s", "https://127.0.0.1", config.SitePortSSL)
	err := http.ListenAndServeTLS(config.SitePortSSL, "cert.pem", "key.pem", mux)
	log.Fatal(err)
}
