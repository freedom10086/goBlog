package models

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"goweb/conf"
	"math/rand"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS `user` (
        `uid` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `username` varchar(20) NOT NULL UNIQUE ,
        `password` varchar(50) NOT NULL,
        `email` varchar(30) NOT NULL DEFAULT '' ,
        `regtime` datetime NOT NULL,
        `sites` varchar(30) NOT NULL DEFAULT '' ,
        `sex` integer NOT NULL DEFAULT 0 ,
        `author` varchar(30) NOT NULL DEFAULT '' ,
        `tags` varchar(100) NOT NULL DEFAULT '' ,
        `status` tinyint NOT NULL DEFAULT 0 ,
		`views` integer NOT NULL DEFAULT 0 ,
        `replys` integer NOT NULL DEFAULT 0
    ) ENGINE=InnoDB;
    CREATE INDEX `post_cid` ON `post` (`fid`);
    CREATE INDEX `post_views` ON `post` (`views`);
*/

type User struct {
	Uid         int
	Username    string
	Password    string
	Email       string
	Regtime     time.Time
	Sites       string
	Sex         int
	Description string
	Exp         int
	Status      int
	Birthday    time.Time
	Phone       string
	Position    string
}

func AddUser(username, password, email string) error {

	_, err := db.Exec(
		"INSERT INTO `user` (`username`,`password`,`email`,`regtime`) VALUES (?,?,?,?)",
		username, password, email, time.Now())
	return err
}

func DelUser(uid int) error {
	_, err := db.Exec(
		"DELETE FROM `user` WHERE uid = ?",
		uid)
	return err
}

func UpdateUser(uid int, email, sites, description, phone, position string, sex int, birthday time.Time) error {
	_, err := db.Exec(
		"UPDATE `user` SET  `email` = ?, `sites` = ?, `description` = ?, `phone` = ?,`position` = ?,`sex` = ?,`birthday` = ? WHERE `uid` = ?",
		email, sites, description, phone, position, sex, birthday, uid)
	return err
}

func GetUser(uid int) (*User, error) {

	//查询数据
	rows, err := db.Query("SELECT  `username`, `email`, `regtime`, `sites`, `sex`, `description`, `exp`, `status`, `birthday`,`phone`,`position` FROM `user` WHERE `uid` = ?",
		uid)

	if err != nil {
		return nil, err
	}
	user := &User{Uid: uid}
	if rows.Next() {
		var timestr string
		var timestrbir string
		err = rows.Scan(&user.Username, &user.Email, &timestr, &user.Sites, &user.Sex, &user.Description, &user.Exp, &user.Status, &timestrbir, &user.Phone, &user.Position)
		loc := time.Local
		user.Regtime, err = time.ParseInLocation(time.RFC822, timestr, loc)
		user.Birthday, err = time.ParseInLocation(time.RFC822, timestr, loc)
		return user, err
	}

	err = errors.New("no user")
	return user, err
}

//存入数据库 md5(password)
func Md5_password(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return md5pass
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
		Expires:  time.Now().Add(time.Second * time.Duration(timeout)),
	}

	mac := hmac.New(sha256.New, []byte(conf.SecretKey))
	strdata, _ := json.Marshal(data)
	mac.Write(strdata)
	signature := mac.Sum(nil)

	return base64.URLEncoding.EncodeToString([]byte(string(strdata) + string(signature)))
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
