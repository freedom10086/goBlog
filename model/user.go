package model

import (
	"database/sql"
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

//添加用户
func AddUser(username, password, email string, sex int) (int64, error) {
	md5pass := Md5_encode(password)
	s := "INSERT INTO `user` (`username`, `password`, `email`, `sex`) VALUES ($1, $2, $3, $4)"
	return add(s, username, md5pass, email, sex)
}

//删除用户
func DelUser(uid int) (int64, error) {
	s := "delete from user where uid =$1"
	return del(s, uid)
}

//更新用户
func UpdateUser(uid int, sex int, birthday, phone, description, site string) (int64, error) {
	s := "UPDATE `user` SET `sex` = $1,`birthday` = $2," +
		"`phone` = $3 `description` = $4, `site` = $5, WHERE `uid` = $6"
	return update(s, sex, birthday, phone, description, site, uid)
}

//更改密码
func UpdatePass1(uid int, password string) (int64, error) {
	md5pass := Md5_encode(password)
	s := "UPDATE `user` SET `password` = $1 WHERE `uid` = $2"
	return update(s, md5pass, uid)
}

//修改密码
func UpdatePass2(uid int, oldpass, newpass string) (int64, error) {
	s := "select username from user where uid = $1 and password = $2"
	oldpass = Md5_encode(oldpass)
	var uname string
	if err := queryA2(s, uid, oldpass, &uname); err != nil {
		return -1, err
	}
	newpass = Md5_encode(newpass)
	s = "UPDATE `user` SET `password` = $ WHERE `uid` = $1"
	return update(s, newpass, uid)
}

//根据uid查询用户
func GetUserById(uid int64) (u *User, err error) {
	u = &User{Uid: uid}
	s := "SELECT `username`,`password`,`email`,`status`,`sex`," +
		"`exp`,`birthday`,`phone`,`description`, " +
		"`site`,`posts`,`replys`,`regtime` " +
		"FROM `user` WHERE `uid` = $1"
	err = queryA1(s, uid, &u.Username, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//根据用户名查询用户
func GetUserByName(username string) (u *User, err error) {
	u = &User{Username: username}
	s := "SELECT `uid`,`password`,`email`,`status`,`sex`," +
		"`exp`,`birthday`,`phone`,`description`, " +
		"`site`,`posts`,`replys`,`regtime` " +
		"FROM `user` WHERE `username` = $1"
	err = queryA1(s, username, &u.Uid, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//根据email查询用户
func GetUserByEmail(email string) (u *User, err error) {
	u = &User{Email: email}
	s := "SELECT `uid`,`password`,`username`,`status`,`sex`," +
		"`exp`,`birthday`,`phone`,`description`, " +
		"`site`,`posts`,`replys`,`regtime` " +
		"FROM `user` WHERE `email` = $1"
	err = queryA1(s, email, &u.Uid, &u.Password, &u.Username, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//username 可能为邮件/用户名
//Status //0-正常 1-禁止访问
func GetUserByNameEmail(username, password string) (u *User, err error) {
	password = Md5_encode(password)
	u = &User{Password: password}
	err = db.QueryRow("SELECT `uid`,`username`,`email`,`status`,`sex`," +
		"`exp`,`birthday`,`phone`,`description`, " +
		"`site`,`posts`,`replys`,`regtime` " +
		"FROM `user` WHERE (`email` = $1 OR `username` = $1) AND `password` = $2",
		username, password).Scan(
		&u.Uid, &u.Username, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//获得所有用户
func GetUsers(page, pagesize int) (us []*User, err error) {
	var rows *sql.Rows
	offset := (page - 1) * pagesize
	if rows, err = db.Query(
		"SELECT `uid`,`username`,`email`, `status`," +
			" `sex`, `exp`, `birthday`, `phone`," +
			" `description`,`site`,`posts`,`replys`,`regtime` " +
			"FROM `user` ORDER BY uid DESC LIMIT $1 OFFSET $2",
		pagesize, offset); err != nil {
		return
	}
	defer rows.Close()
	us = make([]*User, 0, pagesize)

	for rows.Next() {
		user := &User{}
		err = rows.Scan(
			&user.Uid, &user.Username, &user.Email, &user.Status,
			&user.Sex, &user.Exp, &user.Birthday, &user.Phone,
			&user.Description, &user.Site, &user.Posts, &user.Replys, &user.Regtime)

		if err != nil {
			return
		}
		us = append(us, user)
	}

	err = rows.Err()
	return
}

//验证邮箱和用户名
func CheckEmail(email string) bool {
	s := "select count(*) from user where email = $1"
	var num int = 0
	if err := queryA1(s, email, &num); err != nil || num == 0 {
		return true
	}
	return false
}

func CheckUsername(username string) bool {
	s := "select count(*) from user where username = $1"
	var num int = 0
	if err := queryA1(s, username, &num); err != nil || num == 0 {
		return true
	}
	return false
}

//禁止用户
func BlockUser(uid int) (int64, error) {
	s := "UPDATE `user` SET `status` = '1' WHERE `uid` = $1"
	return update(s, uid)
}

//允许用户
func OpenUser(uid int) (int64, error) {
	s := "UPDATE `user` SET `status` = '0' WHERE `uid` = $1"
	return update(s, uid)
}
