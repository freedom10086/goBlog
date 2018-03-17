package router

import (
	"bytes"
	"encoding/base64"
	"goBlog/model"
	"log"
	"net/http"
	"strings"
	"time"
	"io/ioutil"
	"fmt"
	"regexp"
	"encoding/json"
)

type OauthHandler struct {
	BaseHandler
}

func (*OauthHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (h *OauthHandler) DoGet(w http.ResponseWriter, r *http.Request) () {
	if code := r.FormValue("code"); strings.HasPrefix(r.FormValue("state"), "qq_login") { //qq登陆
		// get access_token
		url := fmt.Sprintf("https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
			config.QQConnectAppId, config.QQConnectSecret, code, config.QQConnectRedirect)
		resp, err := http.Get(url)
		if err != nil {
			Unauthorized(w, r, err.Error())
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		s := string(body)
		if strings.Contains(s, "error") {
			Unauthorized(w, r, s)
			return
		}

		//access_token=1771C08BE79247523406A1775DE97C4C&expires_in=7776000&refresh_token=B11F81620A2997E26B149B473327FB3D
		m := make(map[string]string)
		for _, v := range strings.Split(s, "&") {
			m[strings.Split(v, "=")[0]] = strings.Split(v, "=")[1]
		}

		var token string
		token, ok := m["access_token"]
		if !ok {
			Unauthorized(w, r, "获取access_token出错")
			return
		}

		fmt.Println("token:", token)

		// get openId
		url = fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", token)
		resp, err = http.Get(url)
		if err != nil {
			Unauthorized(w, r, err.Error())
			return
		}

		body, err = ioutil.ReadAll(resp.Body)
		s = string(body)
		if strings.Contains(s, "error") {
			Unauthorized(w, r, s)
			return
		}

		// callback( {"client_id":"101462035","openid":"293ABB49EC26DE30AD105E46E2AA051F"} );
		reg := regexp.MustCompile(`"openid":"([0-9A-Z]+)"`)
		result := reg.FindStringSubmatch(s) //[0] 整个字符串 [1] - ()里面的
		if result == nil {
			Unauthorized(w, r, s)
			return
		}

		openId := result[1]
		fmt.Println("openId:", openId)

		// get user info
		url = fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s", token, config.QQConnectAppId, openId)
		resp, err = http.Get(url)
		if err != nil {
			Unauthorized(w, r, err.Error())
			return
		}

		data := &model.QQConnectUserInfoResult{}
		err = json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			InternalError(w, r, err)
			return
		}

		if data.Ret != 0 {
			Unauthorized(w, r, data.Msg)
			return
		}

		http.Redirect(w, r, config.QQConnectRedirect+"?nickname="+data.Nickname+"&access_token="+token+"&state=qq_login_xdluoyang", http.StatusTemporaryRedirect)
		return
	}

	Unauthorized(w, r, "非法的请求")
}

//用来产生token,需要[用户名:密码]base64编码
func (h *OauthHandler) DoPost(w http.ResponseWriter, r *http.Request) {
	if auth := r.Header.Get("Authorization"); auth != "" && strings.Index(auth, "Basic ") == 0 {
		var decodeAuth []byte
		var err error
		if decodeAuth, err = base64.URLEncoding.DecodeString(auth[6:]); err != nil {
			log.Println(err)
			BadParameter(w, r, err.Error())
			return
		}

		if index := bytes.IndexByte(decodeAuth, ':'); index <= 0 {
			log.Println(decodeAuth)
			Unauthorized(w, r, "")
			return
		} else {
			userName := string(decodeAuth[:index])
			passWord := string(decodeAuth[index+1:])
			log.Printf("username is %s password is %s", userName, passWord)

			if u, err := model.UserLogin(userName, passWord); err != nil {
				Unauthorized(w, r, err.Error())
				return
			} else {
				if t, err := model.GenToken(u, 1, config.SecretKey, time.Hour*24*7); err != nil {
					Unauthorized(w, r, err.Error())
					return
				} else {
					Result(w, r, t)
					return
				}
			}
		}
	}

	Unauthorized(w, r, "")
}
