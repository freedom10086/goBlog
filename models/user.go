package models

import (
	"crypto/md5"
	"fmt"
	"log"
	"time"
)

type User struct {
	Uid         int
	Username    string
	Password    string
	Email       string
	Sex         int
	Description string
	Exp         int
	Sites       string
	Messages    int
	Posts       int
	Replys      int
	Phone       string
	Regtime     time.Time
	Birthday    time.Time
}

//添加用户
func AddUser(username, password, email string, sex int) error {
	md5pass := Md5_password(password)
	res, err := db.Exec("call user_reg(?,?,?,?)",
		username, md5pass, email, sex)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//登陆 //Status //0-正常 1-禁止访问
func UserLogin(username, password string) (*User, error) {
	md5pass := Md5_password(password)
	user := &User{}
	status := 0
	err := db.QueryRow("SELECT `uid`,`username`,`email`,`sex`,`exp`,`messages`,`posts`,`replys`,`description`, `status` FROM `user`"+
		" WHERE (`username`=?  AND `password` = ?) OR (`email` = ? AND `password` = ?)", username, md5pass, username, md5pass).Scan(
		&user.Uid, &user.Username, &user.Email, &user.Sex, &user.Exp, &user.Messages, &user.Posts, &user.Replys, &user.Description, &status)

	switch {
	case err != nil:
		log.Fatal(err)
		return nil, err
	case status == 0:
		user, err = GetUserById(user.Uid)
		return user, err
	default:
		return nil, ErrLogin
	}
}

//存入数据库 md5(password)
func Md5_password(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return md5pass
}

//删除用户
func DelUser(uid int) error {
	res, err := db.Exec("call user_del(?)", uid)
	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//更新用户
func UpdateUser(uid int, sites, description, phone string, sex int, birthday time.Time) error {
	res, err := db.Exec("call user_edit(?,?,?,?,?,?)",
		uid, sex, description, sites, birthday, phone)
	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//禁止用户
func BlockUser(uid int) error {
	res, err := db.Exec("call user_bolck(?)", uid)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//允许用户
func OpenUser(uid int) error {
	res, err := db.Exec("call user_open(?)", uid)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}

//修改密码
func ChangePass(uid int, oldpass, newpass string) error {
	uidout := -1
	status := 0
	err := db.QueryRow("SELECT `uid`,`username`, `status` FROM `user` WHERE `uid`=?  AND `password` = ?", uid, oldpass).Scan(
		&uidout, &status)
	if err != nil {
		return err
	}

	if status != 0 {
		return ErrLogin
	}

	res, err := db.Exec("call user_changepass(?,?)", uid, newpass)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return ErrNoAff
	}

	return err
}

//根据uid查询用户
func GetUserById(uid int) (*User, error) {

	auser := &User{Uid: uid}
	err := db.QueryRow("SELECT `username`,`email`,`sex`,`exp`,`messages`,`posts`,`replys`,`description`, `sites`,`phone`,`regtime` "+
		"FROM `user` WHERE `uid` = ?", uid).Scan(
		&auser.Username, &auser.Email, &auser.Sex, &auser.Exp, &auser.Messages, &auser.Posts, &auser.Replys,
		&auser.Description, &auser.Sites, &auser.Phone, &auser.Replys)

	if err != nil {
		return nil, err
	}
	//user.Regtime, err = time.ParseInLocation(time.RFC822, timestr, loc)
	return auser, nil
}

//根据uid查询用户
func GetUserByName(username string) (*User, error) {

	auser := &User{Username: username}
	err := db.QueryRow("SELECT `uid`,`email`,`sex`,`exp`,`messages`,`posts`,`replys`,`description`, `sites`,`phone`,`regtime` "+
		"FROM `user` WHERE `username` = ?", username).Scan(
		&auser.Uid, &auser.Email, &auser.Sex, &auser.Exp, &auser.Messages, &auser.Posts, &auser.Replys,
		&auser.Description, &auser.Sites, &auser.Phone, &auser.Replys)

	if err != nil {
		return nil, err
	}
	//user.Regtime, err = time.ParseInLocation(time.RFC822, timestr, loc)
	return auser, nil
}

//获得所有用户
func GetUsers(order bool, limit, offset int) ([]*User, error) {
	orderstr := "ASC"
	if !order {
		orderstr = "DESC"
	}
	rows, err := db.Query(
		"SELECT `uid`,`username`,`email`, `regtime`,"+
			" `sites`, `sex`, `description`, `exp`,"+
			" `birthday`,`phone`,`posts`,`replys` FROM `user`"+
			" ORDER BY `uid` ? LIMIT ? OFFSET ? ",
		orderstr, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, limit)

	for rows.Next() {
		user := &User{}
		err = rows.Scan(
			&user.Uid, &user.Username, &user.Email, &user.Regtime,
			&user.Sites, &user.Sex, &user.Description, &user.Exp,
			&user.Birthday, &user.Phone, &user.Posts,&user.Replys)

		if err != nil {
			log.Fatal(err)
			continue
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users, err
}


//获得所有被禁止的user
func GetBlockUsers(limit, offset int)([]*User, error){
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
			&user.Sites, &user.Sex, &user.Description, &user.Exp,
			&user.Posts,&user.Replys)

		if err != nil {
			log.Fatal(err)
			continue
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return users, err
}
