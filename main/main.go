package main

import (
	"fmt"
	"goweb/conf"
	"goweb/models"
	"goweb/router"
	"io"
	"log"
	"net/http"
	//"net/url"
)

func init() {
	models.InitDB()
}

func main() {

	fmt.Println("== server start", conf.SiteAddr+conf.SitePort, "==")
	err := models.ModifyCate(7, "cat1mol7", "des1mol7")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	for _, route := range router.Routes {
		mux.HandleFunc(route.Path, route.Func)
	}

	//mux.Handle("/", &myHandler{})
	log.Fatal(http.ListenAndServe(conf.SitePort, mux))
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}
