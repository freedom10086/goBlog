package router

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goBlog/repository"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// GitHub登陆文档
// https://developer.github.com/apps/building-oauth-apps/authorizing-oauth-apps/

type OauthHandler struct {
	BaseHandler
}

func (*OauthHandler) DoAuth(method int, r *http.Request) error {
	return nil
}

func (h *OauthHandler) DoGet(w http.ResponseWriter, r *http.Request) {
	if code := r.FormValue("code"); strings.HasPrefix(r.FormValue("state"), "qq_login") { //qq登陆
		// get access_token
		getUrl := fmt.Sprintf("https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s",
			config.QQConnectAppId, config.QQConnectSecret, code, config.QQConnectRedirect)
		resp, err := http.Get(getUrl)
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

		// get openId and unionid
		// 开发者可通过openID来获取用户的基本信息。
		// 特别需要注意的是，如果开发者拥有多个移动应用、网站应用，可通过获取用户的unionID来区分用户的唯一性，因为只要是同一QQ互联平台下的不同应用，unionID是相同的。
		// 换句话说，同一用户，对同一个QQ互联平台下的不同应用，unionID是相同的
		getUrl = fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s&unionid=1", token)
		resp, err = http.Get(getUrl)
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

		// callback( {"client_id":"101462035","openid":"293ABB49EC26DE30AD105E46E2AA051F","unionid":""} );
		reg := regexp.MustCompile(`"openid":"([0-9A-Z]+)"`)
		result := reg.FindStringSubmatch(s) //[0] 整个字符串 [1] - ()里面的
		if result == nil {
			Unauthorized(w, r, s)
			return
		}

		openId := result[1]
		fmt.Println("openId:", openId)

		// get unionid
		reg = regexp.MustCompile(`"unionid":"([0-9A-Z]+)"`)
		result = reg.FindStringSubmatch(s) //[0] 整个字符串 [1] - ()里面的
		if result == nil {
			Unauthorized(w, r, s)
			return
		}

		unionid := result[1]
		fmt.Println("unionid:", unionid)

		// get user info
		getUrl = fmt.Sprintf("https://graph.qq.com/user/get_user_info?access_token=%s&oauth_consumer_key=%s&openid=%s", token, config.QQConnectAppId, openId)
		resp, err = http.Get(getUrl)
		if err != nil {
			Unauthorized(w, r, err.Error())
			return
		}

		data := &repository.QQConnectUserInfoResult{}
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
	} else if code := r.FormValue("code"); strings.HasPrefix(r.FormValue("state"), "github_login") { // GitHub登陆
		// get code exchange access token
		postUrl := "https://github.com/login/oauth/access_token"

		resp, err := http.PostForm(postUrl, url.Values{
			"client_id":     {config.GitHubClientId},
			"client_secret": {config.GitHubClientSecret},
			"code":          {code},
			"state":         {r.FormValue("state")},
		})
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

		//access_token=e72e16c7e42f292c6912e7710c838347ae178b4a&token_type=bearer
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

		// Use the access token to access the API
		req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
		if err != nil {
			Unauthorized(w, r, err.Error())
			return
		}
		req.Header.Add("Authorization", fmt.Sprintf("token %s", token))
		client := &http.Client{}
		resp, err = client.Do(req)
		if err != nil {
			Unauthorized(w, r, err.Error())
			return
		}
		defer resp.Body.Close()

		data := &repository.GitHubUserInfoResult{}
		err = json.NewDecoder(resp.Body).Decode(data)
		if err != nil {
			InternalError(w, r, err)
			return
		}

		w.Write([]byte("hello " + data.Nickname))
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

			if u, err := repository.UserLogin(userName, passWord); err != nil {
				Unauthorized(w, r, err.Error())
				return
			} else {
				if t, err := repository.GenToken(u, 1, config.SecretKey, time.Hour*24*7); err != nil {
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
