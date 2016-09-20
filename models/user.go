package models

import (
	"crypto/md5"
	"errors"
	"fmt"
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

func Login(username, password string) (bool, error) {
	md5pass := Md5_password(password)

	var uid int
	err := db.QueryRow("SELECT `uid` FROM `user` WHERE `username`=? AND `password` = ?",
		username, md5pass).Scan(&uid)
	if err != nil {
		return false, err
	} else {
		fmt.Println(uid)
		return true, nil
	}
}

//存入数据库 md5(password)
func Md5_password(password string) string {
	md5pass := fmt.Sprintf("%x", md5.Sum([]byte(password)))

	return md5pass
}

func AddUser(username, password, email string) error {
	md5pass := Md5_password(password)
	_, err := db.Exec(
		"INSERT INTO `user` (`username`,`password`,`email`,`regtime`) VALUES (?,?,?,?)",
		username, md5pass, email, time.Now())
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
