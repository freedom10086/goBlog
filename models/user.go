package models

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"goBlog/code"
	"log"
	"math/rand"
	"net/smtp"
	"strings"
	"time"
)

type User struct {
	Uid         int64
	Username    string
	Password    string `json:"-"`
	Email       string
	Status      int
	Sex         int
	Exp         int
	Birthday    time.Time
	Phone       string
	Description string `json:",omitempty"`
	Site        string
	Posts       int
	Replys      int
	Regtime     time.Time
}

type LoginToken struct {
	Uid     int64
	Salt    string
	Expires time.Time
}

type RegToken struct {
	Username string
	Password string //真正要插入时才有此
	Email    string
	Sex      int //同上
	Salt     string
	Expires  time.Time
}

//添加用户
func AddUser(username, password, email string, sex int) (int64, error) {
	md5pass := Md5_encode(password)
	stmt, err := db.Prepare("INSERT INTO `user`(`username`,`password`,`email`,`sex`) VALUES (?,?,?,?)")
	if err != nil {
		return -1, err
	}
	return add(stmt, username, md5pass, email, sex)
}

//删除用户
func DelUser(uid int) (int64, error) {
	stmt, err := db.Prepare("delete from `user` where `uid` = ?")
	if err != nil {
		return -1, err
	}
	return delete(stmt, uid)
}

//更新用户
func UpdateUser(uid int, sites, description, phone string, sex int, birthday time.Time) error {
	stmt, err := db.Prepare("UPDATE `user` SET `sex` = ?,`description` = ?,`sites` = ?,`birthday` = ?,`phone` = ? WHERE `uid` = ?")
	if err != nil {
		return err
	}
	return update(stmt, sex, description, sites, birthday, phone, uid)
}

//获得所有用户
func GetUsers(order bool, offset, limit int) ([]*User, error) {
	//orderstr := "ASC"
	if !order {
		//orderstr = "DESC"
	}
	rows, err := db.Query(
		"SELECT `uid`,`username`,`email`, `status`,"+
			" `sex`, `exp`, `birthday`, `phone`,"+
			" `description`,`site`,`posts`,`replys`,`regtime` "+
			"FROM `user` ORDER BY uid LIMIT ? OFFSET ? ",
		limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, 0)

	for rows.Next() {
		user := &User{}
		err = rows.Scan(
			&user.Uid, &user.Username, &user.Email, &user.Status,
			&user.Sex, &user.Exp, &user.Birthday, &user.Phone,
			&user.Description, &user.Site, &user.Posts, &user.Replys, &user.Regtime)

		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, user)
	}

	err = rows.Err()
	return users, err
}

//根据uid查询用户
func GetUserById(uid int64) (u *User, err error) {
	u = &User{Uid: uid}
	err = db.QueryRow("SELECT `username`,`password`,`email`,`status`,`sex`,"+
		"`exp`,`birthday`,`phone`,`description`, "+
		"`site`,`posts`,`replys`,`regtime` "+
		"FROM `user` WHERE `uid` = ?", uid).Scan(
		&u.Username, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//根据用户名查询用户
func GetUserByName(username string) (u *User, err error) {
	u = &User{Username: username}
	err = db.QueryRow("SELECT `uid`,`password`,`email`,`status`,`sex`,"+
		"`exp`,`birthday`,`phone`,`description`, "+
		"`site`,`posts`,`replys`,`regtime` "+
		"FROM `user` WHERE `username` = ?", username).Scan(
		&u.Uid, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//根据email查询用户
func GetUserByEmail(email string) (u *User, err error) {
	u = &User{Email: email}
	err = db.QueryRow("SELECT `uid`,`password`,`username`,`status`,`sex`,"+
		"`exp`,`birthday`,`phone`,`description`, "+
		"`site`,`posts`,`replys`,`regtime` "+
		"FROM `user` WHERE `email` = ?", email).Scan(
		&u.Uid, &u.Password, &u.Username, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//登陆 //Status //0-正常 1-禁止访问
func UserLogin(username, email, password string) (u *User, err error) {
	if username != "" {
		u, err = GetUserByName(username)
		if err == nil && u != nil && u.Username == username {
			if u.Password == Md5_encode(password) {
				return
			}
		}
	} else if email != "" {
		u, err = GetUserByEmail(email)
		if err == nil && u != nil && u.Email == email {
			if u.Password == Md5_encode(password) {
				return
			}

		}
	}

	if err == nil {
		err = code.ERR_LOGIN
	}

	return nil, err
}

//禁止用户
func BlockUser(uid int) error {
	stmt, err := db.Prepare("UPDATE `user` SET `status` = '1' WHERE `uid` = ?")
	if err != nil {
		return err
	}
	return update(stmt, uid)
}

//允许用户
func OpenUser(uid int) error {
	stmt, err := db.Prepare("UPDATE `user` SET `status` = '0' WHERE `uid` = ?")
	if err != nil {
		return err
	}
	return update(stmt, uid)
}

//修改密码
func ChangePass(uid int, oldpass, newpass string) error {
	status := -1
	err := db.QueryRow("SELECT `status` FROM `user` WHERE `uid` = ? AND `password` = ?", uid, oldpass).Scan(&status)
	switch {
	case err != nil:
		return err
	case status == 0:
		stmt, err := db.Prepare("UPDATE `user` SET `password` = ? WHERE `uid` = ?")
		if err != nil {
			return err
		}
		return update(stmt, newpass, uid)
	default:
		return code.ErrLogin
	}
}

//获得所有被禁止的user
func GetBlockUsers(limit, offset int) ([]*User, error) {
	rows, err := db.Query(
		"SELECT `uid`,`username`,`email`, `regtime`,"+
			" `sites`, `sex`, `description`, `exp`,"+
			" `posts`,`replys` FROM `user`"+
			" WHERE `status` <> 0  LIMIT ? OFFSET ? ",
		limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, limit)

	for rows.Next() {
		user := &User{}
		err = rows.Scan(
			&user.Uid, &user.Username, &user.Email, &user.Regtime,
			&user.Site, &user.Sex, &user.Description, &user.Exp,
			&user.Posts, &user.Replys)

		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, user)
	}

	err = rows.Err()
	return users, err
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
		payload := decode_token[:totallen-32]
		signature := decode_token[totallen-32:]

		mac := hmac.New(sha256.New, []byte(secretKey))
		mac.Write([]byte(payload))
		my_signature := mac.Sum(nil)

		if !hmac.Equal(signature, my_signature) {
			return nil, code.ERR_TOKEN_INVALID
		} else {
			return payload, nil
		}
	}
}

//生成TOKEN base64(data+hmac(data,SecretKey))
func GenLoginToken(uid int64, secretKey string, duration time.Duration) string {
	data := &LoginToken{
		Uid:     uid,
		Salt:    Krand(10),
		Expires: time.Now().Add(duration),
	}
	return EncodeToken(data, secretKey)
}

//valid token 返回uid
func ValidLoginToken(token, secretKey string) (int64, error) {
	if s, err := DecodeToken(token, secretKey); err != nil {
		return -1, nil
	} else {
		data := &LoginToken{}
		if err = json.Unmarshal(s, &data); err != nil {
			return -1, err
		} else if data.Expires.Before(time.Now()) {
			return -1, code.ERR_TOKEN_TIMEOUT
		} else {
			return data.Uid, nil
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
func ValidRegToken(token, secretKey string) (*RegToken, bool) {
	if s, err := DecodeToken(token, secretKey); err != nil {
		log.Println(err)
		return nil, false
	} else {
		data := &RegToken{}
		if err = json.Unmarshal(s, &data); err != nil {
			return nil, false
		} else if data.Expires.After(time.Now()) {
			return data, true
		}
	}

	return nil, false
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

//发送邮箱
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

//发送验证邮箱
func SendValidMail(to, username string) error {
	//sub := "valid your email"
	//timeout := 20
	//token := GenToken(username, timeout)
	//content := "welcome regiest " + configure.SiteName + "!!" +
	//	"\r\n click " + configure.SiteAddr + configure.SitePort + "/email?token=" + token + " to valid your email<" + to + ">" +
	//	"\r\n\r\n attention: please valided in " + string(timeout) + " minutes."
	//return SendToMail(to, sub, content)

	return nil
}

//存入数据库 md5(password)
func Md5_encode(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))
	return md5pass
}
