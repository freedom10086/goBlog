package models

import (
	"errors"
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
