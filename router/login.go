package router

import (
	"database/sql"
	"errors"
	"fmt"
	"goBlog/model"
	"net/http"
	"strconv"
	"time"
)

// todo 线程不安全 use sync.Map
var uuids = make(map[string]client, 500)

const timeout = time.Second * 60

type client struct {
	ch   chan string //往这个里面写入uid表示扫码成功
	time time.Time
}

type LoginHandler struct {
	BaseHandler
}

//二维码登录
type QrLoginHandler struct {
	BaseHandler
}

//登陆返回结果
type LoginResult struct {
	model.User `json:"user"`
	Token      string `json:"token"`
}

func (h *LoginHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

//token null /regiest ->登陆页面
//登陆页面 ->dopost -> 发邮件 -> 点击连接 -> user.doPost 插入数据库
func (h *LoginHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	Template(w, &TemplateData{
		Title: "登陆",
		Css:   []string{"style.css"},
		Js:    []string{"base.js", "particles.js", "qrcode.js"}},
		"page.gohtml", "login.gohtml",
	)
}

//登陆
func (h *LoginHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")
	if username == "" {
		BadParameter(w, r, "用户名密码不能为空")
		return
	}

	user, err := model.UserLogin(username, password)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("此用户不存在")
		}
		InternalError(w, r, err)
		return
	}

	if token, err := model.GenToken(user, 1, config.SecretKey, time.Hour*24*7); err != nil {
		Unauthorized(w, r, err.Error())
		return
	} else {
		Result(w, r, &LoginResult{
			*user,
			token,
		})
	}
}

func (h *QrLoginHandler) DoAuth(method int, r *http.Request) error {
	if method == MethodPost {
		//return h.BaseHandler.DoAuth(method, r)
		return nil
	} else {
		return nil
	}
}

//利用html5 Server-Sent推送
func (h *QrLoginHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Accept") != "text/event-stream" {
		w.Write([]byte("only support text/event-stream"))
		return
	}
	w.Header().Set("Content-Type", "text/event-stream;charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")

	//event:事件类型. uuid连接上 data扫码成功 timeout服务器关闭连接
	//data:消息的数据字段.
	//id:事件ID.
	//retry:一个整数值,指定了重新连接的时间(单位为毫秒),如果该字段值不是整数,则会被忽略.
	//每个字段以\n\n结尾data要换行用\r\n

	//生成不重复的uuid
	uuid := model.GenGuid()
	var ch chan string //0:uid
	for {
		if _, ok := uuids[uuid]; ok {
			uuid = model.GenGuid()
		} else {
			ch = make(chan string, 1)
			uuids[uuid] = client{
				ch,
				time.Now(),
			}
			break
		}
	}

	//发送uuid到客户端
	str := "event: uuid\n"
	str = str + "data: " + uuid + "\n\n"
	str = str + "retry:" + "1500" + "\n\n"
	w.Write([]byte(str))
	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}

	// 遍历uuids清理过期的uuid
	for k, v := range uuids {
		if v.time.Add(timeout).Before(time.Now()) {
			delete(uuids, k)
		}
	}

	//等待1.5s
	time.Sleep(1500 * time.Millisecond)
	expired := make(chan bool, 1)
	go func() { //30s timeout
		time.Sleep(timeout)
		expired <- true
	}()
	//time.AfterFunc(time.Second *30, func() {
	//	timeout <- true
	//})

	var value string //true 确认 false扫码
	for {
		select {
		case value = <-ch: //客户端已经扫码
			str := "event: data\n"
			str = str + "data: " + value + "\n\n"
			_, err := w.Write([]byte(str))
			if err != nil {
				break
			}

			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}

			if value[0] == '0' { //客户端已经扫码
				fmt.Println("已经扫码")
			} else { //客户端已经确认
				fmt.Println("已经确认")
				delete(uuids, uuid)
				return
			}
		case <-expired:
			//timeout s没有回应
			str := "event: timeout\n"
			str = str + "data: close\n\n"
			delete(uuids, uuid)
			_, err := w.Write([]byte(str))
			if err != nil {
				break
			}
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			return
		}
	}
}

//手机扫码post
func (h *QrLoginHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	uuid := r.FormValue("uuid")
	mod := r.FormValue("mod")
	if uuid == "" {
		BadParameter(w, r, "需要uuid参数")
		return
	}

	if mod == "" {
		mod = "normal" //扫码 confirm //确认
	} else if mod != "confirm" && mod != "normal" {
		BadParameter(w, r, "非法的请求参数")
		return
	}

	if v, ok := uuids[uuid]; ok {
		if v.time.Add(timeout).Before(time.Now()) {
			//过期
			Error(w, "二维码已经过期请刷新重试", http.StatusBadRequest)
			return
		}

		//todo token 取得
		var uid int64 = 123
		if mod == "confirm" {
			//扫码确认成功
			v.ch <- "1:" + strconv.FormatInt(uid, 10)
		} else { ////0--扫码成功
			v.ch <- "0:" + strconv.FormatInt(uid, 10)
			v.time = time.Now() //重新计时
			fmt.Println("write 0")
		}

		Result(w, r, true)
		return
	} else {
		Error(w, "二维码无效", http.StatusBadRequest)
		return
	}
}
