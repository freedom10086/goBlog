package main

import (
	"goBlog/handlers"
	"goBlog/models"
	"log"
	"net/http"
	"time"
	"goBlog/router"
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

	r := router.NewRouter();
	server := &http.Server{
		Addr: config.SitePortSSL,
		Handler: r,
		WriteTimeout: 8 * time.Second,
		ReadTimeout:  8 * time.Second,
	}

	log.Printf("https listen on %s%s", "https://127.0.0.1", config.SitePortSSL)
	err := server.ListenAndServeTLS("cert.pem", "key.pem")
	log.Fatal(err)
}
