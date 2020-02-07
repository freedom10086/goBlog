package repository

import (
	"database/sql"
	"errors"
	"time"
)

type User struct {
	Id          int64      `json:"id"`
	Username    string     `json:"username"`
	Password    string     `json:"-"`
	Email       string     `json:"email"`
	Status      int        `json:"status"`
	Sex         int        `json:"sex"`
	Exp         int        `json:"exp"`
	Birthday    NullTime   `json:"birthday,omitempty"`
	Phone       NullString `json:"phone,omitempty"`
	Description NullString `json:"description,omitempty"`
	Site        NullString `json:"site,omitempty"`
	Posts       int        `json:"posts"`
	Replys      int        `json:"replys"`
	Regtime     time.Time  `json:"regtime"`
}

//添加用户
func AddUser(username, password, email string, sex int) (int, error) {
	md5pass := Md5_encode(password)
	s := "INSERT INTO users (username, password, email, sex) VALUES ($1, $2, $3, $4) RETURNING id"
	return add(s, username, md5pass, email, sex)
}

//删除用户
func DelUser(uid int) (int64, error) {
	s := "delete from users where id =$1"
	return del(s, uid)
}

//更新用户
func UpdateUser(uid int, sex int, birthday, phone, description, site string) (int64, error) {
	s := "UPDATE users SET sex = $1,birthday = $2," +
		"phone = $3 description = $4, site = $5, WHERE id = $6"
	return update(s, sex, birthday, phone, description, site, uid)
}

//更改密码
func UpdatePass1(uid int, password string) (int64, error) {
	md5pass := Md5_encode(password)
	s := "UPDATE users SET password = $1 WHERE id = $2"
	return update(s, md5pass, uid)
}

//修改密码
func UpdatePass2(uid int, oldpass, newpass string) (int64, error) {
	s := "SELECT username FROM users WHERE id = $1 AND password = $2"
	oldpass = Md5_encode(oldpass)
	var uname string
	if err := db.QueryRow(s, uid, oldpass).Scan(&uname); err != nil {
		return -1, err
	}
	newpass = Md5_encode(newpass)
	s = "UPDATE users SET password = $ WHERE id = $1"
	return update(s, newpass, uid)
}

//根据uid查询用户
func GetUserById(uid int64) (u *User, err error) {
	u = &User{Id: uid}
	s := `SELECT username,password,email,status,sex,exp,birthday,phone,description,
		site,posts,replys,created FROM users WHERE id = $1`
	err = db.QueryRow(s, uid).Scan(&u.Username, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//根据用户名查询用户
func GetUserByName(username string) (u *User, err error) {
	u = &User{Username: username}
	s := `SELECT id,password,email,status,sex,exp,birthday,phone,description,
		site,posts,replys,created FROM users WHERE username = $1`
	err = db.QueryRow(s, username).Scan(&u.Id, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//根据email查询用户
func GetUserByEmail(email string) (u *User, err error) {
	u = &User{Email: email}
	s := `SELECT id,password,username,status,sex,exp,birthday,phone,description,
		site,posts,replys,created FROM users WHERE email = $1`
	err = db.QueryRow(s, email).Scan(&u.Id, &u.Password, &u.Username, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)
	return
}

//username 可能为邮件/用户名
//Status //0-正常 1-禁止访问
func UserLogin(username, password string) (u *User, err error) {
	password = Md5_encode(password)
	u = &User{ /*Password: password*/ }
	s := `SELECT id,username,password,email,status,sex,exp,birthday,
	phone,description,site,posts,replys,created FROM users
		WHERE (email = $1 OR username = $1)` //AND password = $2

	err = db.QueryRow(s, username).Scan(
		&u.Id, &u.Username, &u.Password, &u.Email, &u.Status, &u.Sex,
		&u.Exp, &u.Birthday, &u.Phone, &u.Description,
		&u.Site, &u.Posts, &u.Replys, &u.Regtime)

	if err == nil {
		if password != u.Password {
			err = errors.New("密码错误")
			return
		}

		if u.Status < 0 {
			err = errors.New("你已经被封禁，请联系管理员解封")
			return
		}
	}

	return
}

//获得所有用户
func GetUsers(page, pagesize int) (us []*User, err error) {
	var rows *sql.Rows
	offset := (page - 1) * pagesize
	s := `SELECT id,username,email, status,sex, exp, birthday, phone,
		description,site,posts,replys,created FROM users
		ORDER BY id DESC LIMIT $1 OFFSET $2`
	if rows, err = db.Query(s, pagesize, offset); err != nil {
		return
	}
	defer rows.Close()
	us = make([]*User, 0, pagesize)

	for rows.Next() {
		user := &User{}
		err = rows.Scan(
			&user.Id, &user.Username, &user.Email, &user.Status,
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
	s := "SELECT id FROM users WHERE email = $1"
	var num int = 0
	if err := db.QueryRow(s, email).Scan(&num); err == sql.ErrNoRows || num == 0 {
		return true
	}
	return false
}

//check用户名 true ->可用
func CheckUsername(username string) bool {
	s := "SELECT id FROM users WHERE username = $1"
	var num int = 0
	if err := db.QueryRow(s, username).Scan(&num); err == sql.ErrNoRows || num == 0 {
		return true
	}
	return false
}

//禁止用户
func BlockUser(uid int) (int64, error) {
	s := "UPDATE users SET status = '1' WHERE id = $1"
	return update(s, uid)
}

//允许用户
func OpenUser(uid int) (int64, error) {
	s := "UPDATE users SET status = '0' WHERE id = $1"
	return update(s, uid)
}
