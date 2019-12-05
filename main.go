package main

import (
	"goBlog/conf"
	"goBlog/logger"
	"goBlog/model"
	"goBlog/router"
	"log"
	"net/http"
	"time"
)

//配置文件
var config *conf.Config

//路由列表
var routers map[string]router.Handler

func init() {
	config = conf.Conf

	//todo
	model.InitDB(config.DbHost, config.DbPort, config.DbName, config.DbUsername, config.DbPassword)
	logger.I("==%s started==", config.SiteName)

	routers = map[string]router.Handler{
		"/":          &router.HomeHandler{},
		"/article":   &router.ArticleHandler{},
		"/categorys": &router.CateHandler{},
		"/users":     &router.UserHandler{},
		"/oauth":     &router.OauthHandler{},
		"/login":     &router.LoginHandler{},
		"/qrlogin":   &router.QrLoginHandler{},
		"/register":  &router.RegisterHandler{},
		"/chats":     &router.ChatHandler{},
		"/account/":  &router.AccountHandler{}, //reset_password reg_compete...
		"/admin/":    &router.AdminHandler{},
		"/location":  &router.LocationHandler{},
	}
}

func main() {
	defer model.CloseDB()
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			path := "https://127.0.0.1" + config.SitePortSSL + r.URL.Path
			router.Redirect(w, r, path, http.StatusMovedPermanently)
		})
		err := http.ListenAndServe(config.SitePort, nil)
		if err != nil {
			log.Fatal("start server error!", err)
		}
	}()

	r := router.NewRouter()
	r.DefaultHandler(&router.DefaultHandler{})
	for key, value := range routers {
		r.Register(key, value)
	}

	server := &http.Server{
		Addr:           config.SitePortSSL,
		Handler:        r,
		WriteTimeout:   8 * time.Second,
		ReadTimeout:    8 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("https listen on %s%s", "https://127.0.0.1", config.SitePortSSL)
	log.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}
