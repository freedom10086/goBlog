package models

import (
	"log"
	"time"
)

//单一post
type Post struct {
	Tid       int //id
	Cid       int
	User      User
	Title     string
	Content   string
	Type      int
	Status    int
	Views     int
	Replys    int
	Created   time.Time
	Updated   time.Time
	Lastreply time.Time
}

//带回复
type Article struct {
	Post     Post //id
	Comments []*Comment
}

//查看一篇文章是否可以回复
func PostCanReply(tid int) bool {
	status := 2
	err := db.QueryRow(
		"SELECT  `status` FROM `post` WHERE `tid` = ?", tid).Scan(&status)
	if err != nil {
		log.Fatal(err)
		return false
	}

	return status == 0
}

//发布主题
func AddPost(cid, uid int, title, content string) error {
	res, err := db.Exec(
		"call post_add(?,?,?,?)", cid, uid, title, content)

	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return ErrNoAff
	}
	return err
}


//获得文章带回复
func GetArticle(tid int)(*Article,error){
	row := db.QueryRow(
		"SELECT `cid`,`uid`,`author`,`title`, `content`,"+
			"`type`,`status`,`created`, `updated`,"+
			"`views`, `replys`,`lastreply`"+
			" FROM `post` WHERE `tid` = ?",tid)

	post := &Post{Tid: tid}
	err := row.Scan(
		&post.Cid, &post.User.Uid,&post.User.Username, &post.Title,&post.Content,
		&post.Type,&post.Status,&post.Created,&post.Updated,
		&post.Views,&post.Replys,&post.Lastreply)

	//sql.ErrNoRows
	if err != nil {
		return nil, err
	}

	comments,err := GetComments(tid)
	if err != nil {
		log.Fatal(err)
		comments = nil
	}

	artcle :=&Article{Post:post,Comments:comments}
	return artcle,nil
}

//按照创建顺序倒叙
func GetPostNewCreate(cid,limit,offset int) ([]*Post, error){
	return GetPosts(cid, limit, offset,"create")
}

//按照最后回复顺序倒叙
func GetPostNewReply(cid,limit,offset int) ([]*Post, error){
	return GetPosts(cid, limit, offset,"new")
}

//按照最近7天热帖排序
func GetPostHot(cid,limit,offset int) ([]*Post, error){
	return GetPosts(cid, limit, offset,"hot")
}

//获得某一cid的文章列表 按照发布时间倒叙
func GetPosts(cid, limit, offset int,order string) ([]*Post, error) {
	//查询数据
	where  := "ORDER BY `tid` DESC"
	switch order {
	case "create":
		where = "ORDER BY `tid` DESC"
	case "hot":
		//最近7天的热帖
		where = "AND DATEDIFF(NOW(),`lastreply`)< 7 ORDER BY `replys` DESC,`lastreply` DESC"
	case "new":
		//新帖
		where = "ORDER BY `lastreply` DESC"
	default:
		//新帖
		where = "ORDER BY `lastreply` DESC"
	}
	
	rows, err := db.Query(
		"SELECT `tid`,`uid`,`author`,`title`, `content`,"+
			"`type`,`status`,`created`, `updated`,"+
			"`views`, `replys`,`lastreply`"+
			" FROM `post` WHERE `cid` = ? "+where +" LIMIT ? OFFSET ?",
		cid, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, limit)

	for rows.Next() {
		post := &Post{Cid:cid}
		err = rows.Scan(
			&post.Tid, &post.User.Uid,&post.User.Username, &post.Title,&post.Content,
			&post.Type,&post.Status,&post.Created,&post.Updated,
			&post.Views,&post.Replys,&post.Lastreply)

		if err != nil {
			log.Fatal(err)
			continue
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return posts, err
}


//删除主题
func DelPost(id int) error {
	res, err := db.Exec("call post_del(?)",
		id, id)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	log.Println("aff", rowCnt)
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return ErrNoAff
	}
	return err
}



/*
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
