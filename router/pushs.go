package router

//todo
//http/2 push 列表,渲染模板时候要用
//key path,value要push 的资源
var pushs map[string][]string = map[string][]string{
	"/": []string{"styles/base.css"},
}
