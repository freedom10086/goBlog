package model

import (
	"time"
	"encoding/json"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"log"
	"math/rand"
)

var (
	ErrTokenInvalid  = errors.New("token is invalid!")
	ErrTokenExperied = errors.New("token is experied!")
)

type Token struct {
	Uid       int64
	Authority int //权限 0-99 普通用户 100admin
	Salt      string
	Expires   time.Time
}

type RegToken struct {
	Username string
	Email    string
	Salt     string
	Expires  time.Time
}

//生成token
func EncodeToken(data interface{}, secretKey string) string {
	mac := hmac.New(sha256.New, []byte(secretKey))
	bytedata, _ := json.Marshal(data)
	mac.Write(bytedata)
	signature := mac.Sum(nil)
	return base64.URLEncoding.EncodeToString(append(bytedata, signature...))
}

//验证token 签名是否一致
func DecodeToken(token, secretKey string) ([]byte, error) {
	if decode_token, err := base64.URLEncoding.DecodeString(token); err != nil {
		return nil, err
	} else {
		totallen := len(decode_token)
		if totallen < 32 {
			return nil, errors.New("token error")
		}
		payload := decode_token[:totallen-32]
		signature := decode_token[totallen-32:]

		mac := hmac.New(sha256.New, []byte(secretKey))
		mac.Write([]byte(payload))
		my_signature := mac.Sum(nil)

		if !hmac.Equal(signature, my_signature) {
			return nil, ErrTokenInvalid
		} else {
			return payload, nil
		}
	}
}

//生成TOKEN base64(data+hmac(data,SecretKey))
func GenToken(user *User, authority int, secretKey string, duration time.Duration) (string, error) {
	data := &Token{
		Uid:       user.Uid,
		Authority: authority,
		Salt:      Krand(10),
		Expires:   time.Now().Add(duration),
	}
	return EncodeToken(data, secretKey), nil
}

//valid token 返回Token
func ValidToken(token, secretKey string) (*Token, error) {
	if s, err := DecodeToken(token, secretKey); err != nil {
		return nil, err
	} else {
		data := &Token{}
		if err = json.Unmarshal(s, &data); err != nil {
			return nil, err
		} else if data.Expires.Before(time.Now()) {
			return nil, ErrTokenExperied
		} else {
			return data, nil
		}
	}
}

//产生注册token
//包含 username,email,过期时间,产生token发送到邮箱
//当用户点击连接时验证token和时间等，成功则设置密码等完成注册写入数据库
func GenRegToken(username, email string, secretKey string, duration time.Duration) string {
	data := &RegToken{
		Username: username,
		Email:    email,
		Salt:     Krand(10),
		Expires:  time.Now().Add(duration),
	}
	return EncodeToken(data, secretKey)
}

//返回用户名和邮箱以便于检查可用性
func ValidRegToken(token, secretKey string) (*RegToken, error) {
	if s, err := DecodeToken(token, secretKey); err != nil {
		log.Println(err)
		return nil, err
	} else {
		data := &RegToken{}
		if err = json.Unmarshal(s, &data); err != nil {
			return nil, err
		} else if data.Expires.After(time.Now()) {
			return data, nil
		} else {
			return nil, errors.New("token已过期")
		}
	}
}

// 随机字符串
func Krand(size int) string {
	kinds, result := [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		ikind := rand.Intn(3)
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return string(result)
}
