package models

import (
	"time"
)

type Post struct {
	Id       int //id
	Title    string
	Content  string
	Created  time.Time
	Updated  time.Time
	Views    int
	Replys   int
	Status   int //0--open 1--close
	User     User
	Category Category
	Comments []*Comment
}

/*
//发布主题
func AddPost(cid, uid int, title, content, author string) error {
	_, err := db.Exec(
		"INSERT INTO `post` (`cid`,`pid`,`uid`,`title`,`content`,`created`,`author`) VALUES (?,?,?,?,?,?,?)",
		cid, 0, uid, title, content, time.Now(), author)

	if err != nil {
		return nil
	}
	//更新category计数
	_, err = db.Exec("UPDATE `category` SET `posts` = `posts` + 1 WHERE cid = ?",
		cid)
	return err
}

//发表回复 回复对象--id为pid
func AddComment(cid, pid, uid int, content, author string) error {
	//todo 检查status看看是否可以回复

	_, err := db.Exec(
		"INSERT INTO `post` (`cid`,`pid`,`uid`,`content`,`created`,`author`) VALUES (?,?,?,?,?,?)",
		cid, pid, uid, content, time.Now(), author)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"UPDATE `post` SET  `replys` = `replys`+ 1  WHERE `id` = ?",
		pid)
	return err
}

//获得文章根据id
func GetPost(id int) (*Post, error) {

	row := db.QueryRow("SELECT  `cid`, `pid`,`uid`,`title`, `content`, `created`, `updated`, `views`, `replys`, `author`, `status` FROM `post` WHERE `id` = ?",
		id)

	post := &Post{Id: id}
	err := row.Scan(
		&post.Category.Cid, &post.Pid,
		&post.User.Uid, &post.Title,
		&post.Content, &post.Created,
		&post.Updated, &post.Views,
		&post.Replys, &post.User.Username,
		&post.Status)

	//sql.ErrNoRows
	if err != nil {
		return nil, err
	}

	//这是主题 增加阅读量
	if post.Pid == 0 {
		if _, err = db.Exec("UPDATE `post` SET  `views` = `views`+ 1  WHERE `id` = ?", id); err != nil {
			log.Fatal(err)
		}
	}

	//获得回复
	if post.Replys > 0 {
		post.Comments, err = GetComments(id)
		if err != nil {
			log.Fatal(err)
		}
	}
	return post, nil
}

//获得某一文章的所有评论/或者某一评论的子评论
func GetComments(id int) ([]*Comment, error) {

	rows, err := db.Query(
		"SELECT `id`,`uid`,`content`,`created`,`updated`,replys` FROM `post` WHERE `pid` = ?",
		id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := make([]*Comment, 10)
	for rows.Next() {
		comment := &Comment{}
		err = rows.Scan(
			&comment.Id,
			&comment.User.Uid,
			&comment.Content,
			&comment.Created,
			&comment.Updated,
			&comment.Replys)

		if err != nil {
			log.Fatal(err)
		}
		continue
		comments = append(comments, comment)
	}

	return comments, err

}

//获得某一cid的文章列表
func GetPosts(cid, limit, offset int) ([]*Post, error) {
	//查询数据
	rows, err := db.Query(
		"SELECT `id`,`uid`,`title`, `content`,"+
			"`created`, `updated`, `views`, `replys`,"+
			"`author`, `status` FROM `post` WHERE "+
			"`cid` = ? AND `pid` = 0 ORDER BY `id` DESC LIMIT ? OFFSET ? ",
		cid, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, limit)

	for rows.Next() {
		post := &Post{}
		err = rows.Scan(
			&post.Category.Cid, &post.Pid,
			&post.User.Uid, &post.Title,
			&post.Content, &post.Created,
			&post.Updated, &post.Views,
			&post.Replys, &post.User.Username,
			&post.Status)

		if err != nil {
			log.Fatal(err)
			continue
		}

		posts = append(posts, post)
	}

	return posts, err
}

//删除主题
func DelPost(id int) error {
	_, err := db.Exec(
		"DELETE FROM `post` WHERE id = ? or pid = ?",
		id, id)
	if err != nil {
		return err
	}

	//更新category计数
	_, err = db.Exec("UPDATE `category` SET `posts` = `posts` - 1 WHERE cid = ?",
		id)

	return err
}

//删除回复
func DelComment(id, pid int) error {
	err := DelPost(id)
	if err != nil {
		return err
	}

	//更新回复计数
	_, err = db.Exec(
		"UPDATE `post` SET  `replys` = (SELECT COUNT(*) FROM `post` WHERE `pid` = ?)  WHERE `id` = ?",
		pid, pid)

	return err
}

//修改主题或者回复
func ModifyPost(id int, title, content string) error {

	_, err := db.Exec(
		"UPDATE `post` SET  `title` = ?, `content` = ?, `updated` = ? WHERE `id` = ?",
		title, content, time.Now(), id)

	return err
}
*/
