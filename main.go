package main

import (
	"fmt"
	"goweb/models"
	"goweb/router"
	"io"
	"log"
	"net/http"
	"net/url"
)

func init() {
	models.InitDB()
}

func main() {

	//todo 合并评论和主题
	//一些model重写

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

	//fmt.Println(models.Md5_password("justfrsdgfxfdxfsice"))

	//go controls.SendValidMail("1770626192@qq.com", "hehe")

	//err := models.AddUser("hehe2", "password", "2351386755@qq.com")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//models.Login("hehe2", "password")

	/*
		发送post请求


		//这里添加post的body内容
		postUrl := "https://sms.yunpian.com/v2/sms/single_send.json"
		apikey := "744653c93c6355e2dd705a06a6724cdc"
		mobile := "18706798706"
		text := "【出驾学车】您的验证码是6666"

		data := make(url.Values)
		data["apikey"] = []string{apikey}
		data["mobile"] = []string{mobile}
		data["text"] = []string{text}

		//把post表单发送给目标服务器
		res, err := http.PostForm(postUrl, data)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		defer res.Body.Close()

		fmt.Println("post send success")
	*/

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
