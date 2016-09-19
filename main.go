package main

import (
	"fmt"
	"goweb/models"
	"goweb/router"
	"io"
	"log"
	"net/http"
)

func init() {
	models.InitDB()
}

func main() {
	//fmt.Println("=====start server=====")
	//err := models.AddPost(0, 0, "title2", "fuckyou2")

	//models.DelPost(5)
	//post, err := models.GetPost(7)
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	fmt.Println(post)
	//}
	//models.ModifyPost(2, "titlemodify", "contentmodify", "tags,tags2,tags3")
	//posts, err := models.GetPosts(10)
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	for _, post := range posts {
	//		fmt.Println(post)
	//	}
	//}

	//token := models.GenToken("password", 60)
	//fmt.Println(token)

	//isok, message := models.ValidToken(token)

	//fmt.Println(isok, message)

	fmt.Println(models.Md5_password("justfrsdgfxfdxfsice"))

	mux := http.NewServeMux()

	for _, route := range router.Routes {
		mux.HandleFunc(route.Path, route.Func)
	}

	//mux.Handle("/", &myHandler{})
	log.Fatal(http.ListenAndServe(":8080", mux))
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello")
}
