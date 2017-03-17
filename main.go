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
	if err := model.SendMail("2351386755@qq.com", "testemmial", "测试内容"); err != nil {
		log.Println(err)
	}

	defer model.CloseDB()
	go func() {
		http.ListenAndServe(config.SitePort, &router.RedirectHandler{
			Url:  "https://" + config.SiteIpAddr,
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
