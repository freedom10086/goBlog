package models

import (
	"errors"
	"log"
	"time"
)

/*
CREATE TABLE IF NOT EXISTS `post` (
        `tid` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `uid` integer NOT NULL DEFAULT 0 ,
        `fid` integer NOT NULL DEFAULT 0 ,
        `title` varchar(255) NOT NULL DEFAULT '' ,
        `content` varchar(5000) NOT NULL DEFAULT '' ,
        `created` datetime NOT NULL,
        `updated` datetime,
        `author` varchar(30) NOT NULL DEFAULT '' ,
        `tags` varchar(100) NOT NULL DEFAULT '' ,
        `status` tinyint NOT NULL DEFAULT 0 ,
		`views` integer NOT NULL DEFAULT 0 ,
        `replys` integer NOT NULL DEFAULT 0
    ) ENGINE=InnoDB;
    CREATE INDEX `post_cid` ON `post` (`fid`);
    CREATE INDEX `post_views` ON `post` (`views`);
*/
type Post struct {
	Tid     int
	Uid     int
	Fid     int
	Title   string
	Content string
	Created time.Time
	Updated time.Time
	Views   int
	Replys  int
	Author  string
	Tags    string
	Status  int
}

func AddPost(uid, fid int, title, content string) error {
	author := "author" + string(uid)
	timenow := time.Now()
	_, err := db.Exec(
		"INSERT INTO `post` (`uid`,`fid`,`title`,`content`,`created`,`updated`,`author`) VALUES (?,?,?,?,?,?,?)",
		uid, fid, title, content, timenow, timenow, author)
	return err
}

func DelPost(tid int) error {
	_, err := db.Exec(
		"DELETE FROM post WHERE tid = ?",
		tid)
	return err
}

func GetPost(tid int) (*Post, error) {
	//查询数据
	rows, err := db.Query("SELECT  `uid`, `fid`, `title`, `content`, `created`, `updated`, `views`, `replys`, `author`, `tags`, `status` FROM `post` WHERE `tid` = ?",
		tid)

	if err != nil {
		return nil, err
	}
	post := &Post{Tid: tid}

	if rows.Next() {
		var timestr string
		err = rows.Scan(&post.Uid, &post.Fid, &post.Title, &post.Content, &timestr, &post.Updated, &post.Views, &post.Replys, &post.Author, &post.Tags, &post.Status)
		loc := time.Local
		post.Created, err = time.ParseInLocation("2006-01-02 15:04:05", timestr, loc)

		//add view count
		if _, err2 := db.Exec("UPDATE `post` SET  `views` = `views`+ 1  WHERE `tid` = ?", tid); err2 != nil {
			log.Fatal(err2)
		}

		return post, err
	}

	err = errors.New("no post")
	return post, err
}

func GetPosts(limit int) ([]*Post, error) {
	//查询数据
	rows, err := db.Query("SELECT `tid`, `uid`, `fid`, `title`, `content`, `created`, `updated`, `views`, `replys`, `author`, `tags`, `status` FROM `post` LIMIT  ?",
		limit)

	if err != nil {
		return nil, err
	}

	posts := make([]*Post, 0)

	for rows.Next() {
		var timestr string
		var updatestr string
		post := &Post{}
		err = rows.Scan(&post.Tid, &post.Uid, &post.Fid, &post.Title, &post.Content, &timestr, &updatestr, &post.Views, &post.Replys, &post.Author, &post.Tags, &post.Status)
		loc := time.Local

		post.Created, err = time.ParseInLocation("2006-01-02 15:04:05", timestr, loc)
		post.Updated, err = time.ParseInLocation("2006-01-02 15:04:05", updatestr, loc)

		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, post)
	}

	return posts, err
}

func ModifyPost(tid int, title, content, tags string) error {

	_, err := db.Exec(
		"UPDATE `post` SET  `title` = ?, `content` = ?, `updated` = ?, `tags` = ? WHERE `tid` = ?",
		title, content, time.Now(), tags, tid)

	return err
}
