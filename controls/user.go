package controls

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goweb/conf"
	"io"
	"math/rand"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

func UserHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello this is user!!")
}

type TokenData struct {
	Username string
	Salt     string
	Expires  time.Time
}

//生成TOKEN base64(data+hmac(data,SecretKey))
func GenToken(username string, timeout int) string {

	data := &TokenData{
		Username: username,
		Salt:     Krand(10),
		Expires:  time.Now().Add(time.Minute * time.Duration(timeout)),
	}

	mac := hmac.New(sha256.New, []byte(conf.SecretKey))
	bytedata, _ := json.Marshal(data)
	mac.Write(bytedata)
	signature := mac.Sum(nil)

	return base64.URLEncoding.EncodeToString(append(bytedata, signature...))
}

// 随机字符串
func Krand(size int) string {
	kinds, result := [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		ikind := rand.Intn(3)
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}

//valid token
func ValidToken(token string) (bool, string) {

	fmt.Println(len(token))
	decode_token, err := base64.URLEncoding.DecodeString(token)

	if err != nil {
		return false, "decode error"
	}

	totallen := len(decode_token)
	payload := decode_token[:totallen-32]
	signature := decode_token[totallen-32:]

	mac := hmac.New(sha256.New, []byte(conf.SecretKey))
	mac.Write([]byte(payload))
	my_signature := mac.Sum(nil)

	if hmac.Equal(signature, my_signature) {

		data := &TokenData{}
		err := json.Unmarshal([]byte(payload), &data)

		if err != nil {
			return false, "Unmarshal faild"
		}
		if data.Expires.Before(time.Now()) {
			return false, "time is expires"
		}

		return true, data.Username

	} else {
		return false, "signature not equal"
	}
}

func SendToMail(to, subject, content string) error {
	user := "2351386755@qq.com"
	password := "StrikeFreedom"
	host := "smtp.qq.com:587"

	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])

	msg := []byte("To:" + to +
		"\r\n" +
		"Subject: " + subject +
		"\r\n\r\n" +
		content)

	err := smtp.SendMail(host, auth, user, strings.Split(to, ";"), msg)
	return err
}

func SendValidMail(to, username string) error {
	sub := "valid your email"
	timeout := 20
	token := GenToken(username, timeout)
	content := "welcome regiest " + conf.SiteName + "!!" +
		"\r\n click " + conf.Host + "/email?token=" + token + " to valid your email<" + to + ">" +
		"\r\n\r\n attention: please valided in " + string(timeout) + " minutes."
	return SendToMail(to, sub, content)
}
