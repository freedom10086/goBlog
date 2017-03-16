package main

import (
	"goBlog/conf"
	"goBlog/model"
	"goBlog/router"
	"log"
	"net/http"
	"time"
)

var config *conf.Config

func init() {
	config = conf.Conf
	model.InitDB(config.DbName, config.DbUsername, config.DbPassword)
	log.Printf("==%s started==", config.SiteName)
}

func main() {
	defer model.CloseDB()
	go func() {
		log.Printf("http listen on %s%s", config.SiteAddr, config.SitePort)
		http.ListenAndServe(config.SitePort, &router.RedirectHandler{
			Url:  "https://127.0.0.1",
			Port: config.SitePortSSL,
		})
	}()

	r := router.NewRouter()
	server := &http.Server{
		Addr:         config.SitePortSSL,
		Handler:      r,
		WriteTimeout: 8 * time.Second,
		ReadTimeout:  8 * time.Second,
	}

	log.Printf("https listen on %s%s", "https://127.0.0.1", config.SitePortSSL)
	err := server.ListenAndServeTLS("cert.pem", "key.pem")
	log.Fatal(err)
}
