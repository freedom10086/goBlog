package models

import (
	"crypto/md5"
	"errors"
	"fmt"
	"log"
	"time"
)

/*
DROP TABLE IF EXISTS `user`

CREATE TABLE IF NOT EXISTS `user` (
        `uid` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `username` varchar(20) UNIQUE NOT NULL,
        `password` varchar(50) NOT NULL,
        `email` varchar(30) NOT NULL DEFAULT '' ,
        `regtime` datetime NOT NULL,
        `sites` varchar(30) NOT NULL DEFAULT '' ,
        `sex` tinyint NOT NULL DEFAULT 0 ,
        `description` varchar(150) NOT NULL DEFAULT '' ,
        `exp` integer NOT NULL DEFAULT 0 ,
        `status` tinyint NOT NULL DEFAULT 0 ,
		`birthday` date NOT NULL DEFAULT '0000-00-00' ,
        `phone` varchar(20) NOT NULL DEFAULT '',
		`position` varchar(100) NOT NULL DEFAULT ''
    ) ENGINE=InnoDB;
    CREATE INDEX `user_username` ON `user` (`username`);
	CREATE INDEX `user_email` ON `user` (`email`);
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
	Status      int //0--needvalidemail 1--ok  2--block
	Birthday    time.Time
	Phone       string
	Position    string
}

//登陆
func Login(username, password string) (*User, error) {
	md5pass := Md5_password(password)

	user := &User{}

	err := db.QueryRow("SELECT `uid` `status` FROM `user` WHERE `username`=?  AND `password` = ?",
		username, md5pass).Scan(&user.Uid, &user.Status)

	switch {
	case err != nil:
		log.Fatal(err)
		return nil, err
	case user.Status == 0:
		return nil, errors.New("you need valid your email !!")

	case user.Status == 1:
		user, err = GetUser(user.Uid)
		return user, err
	default:
		return nil, errors.New("you dont have permission to access that !!")
	}
}

//存入数据库 md5(password)
func Md5_password(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return md5pass
}

//添加用户
func AddUser(username, password, email string) error {
	md5pass := Md5_password(password)
	_, err := db.Exec(
		"INSERT INTO `user` (`username`,`password`,`email`,`regtime`) VALUES (?,?,?,?)",
		username, md5pass, email, time.Now())
	return err
}

//删除用户
func DelUser(uid int) error {
	_, err := db.Exec(
		"DELETE FROM `user` WHERE uid = ?",
		uid)
	return err
}

//更新用户
func UpdateUser(uid int, email, sites, description, phone, position string, sex int, birthday time.Time) error {
	_, err := db.Exec(
		"UPDATE `user` SET  `email` = ?, `sites` = ?, `description` = ?, `phone` = ?,`position` = ?,`sex` = ?,`birthday` = ? WHERE `uid` = ?",
		email, sites, description, phone, position, sex, birthday, uid)
	return err
}

//根据uid查询用户
func GetUser(uid int) (*User, error) {

	//查询数据
	row := db.QueryRow("SELECT  `username`, `email`, `regtime`, `sites`, `sex`, `description`, `exp`, `status`, `birthday`,`phone`,`position` FROM `user` WHERE `uid` = ?",
		uid)

	var timestr string
	var timestrbir string

	user := &User{Uid: uid}
	err := row.Scan(&user.Username, &user.Email, &timestr, &user.Sites, &user.Sex, &user.Description, &user.Exp, &user.Status, &timestrbir, &user.Phone, &user.Position)

	if err != nil {
		return nil, err
	}
	loc := time.Local
	user.Regtime, err = time.ParseInLocation(time.RFC822, timestr, loc)
	user.Birthday, err = time.ParseInLocation(time.RFC822, timestr, loc)
	return user, err
}

//获得所有用户
func GetUsers(order bool, limit, offset int) ([]*User, error) {

	orderstr := "ASC"
	if !order {
		orderstr = "DESC"
	}

	rows, err := db.Query(
		"SELECT `uid`,`username`,`email`, `regtime`,"+
			"`sites`, `sex`, `description`, `exp`,"+
			"`status`, `birthday`,`phone`,`position` FROM `user` "+
			"ORDER BY `uid` ? LIMIT ? OFFSET ? ",
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
			&user.Status, &user.Birthday, &user.Phone, &user.Position,
		)

		if err != nil {
			log.Fatal(err)
			continue
		}

		users = append(users, user)
	}

	return users, err
}
