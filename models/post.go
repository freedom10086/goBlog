package models

import (
	"log"
	"time"
	"goBlog/code"
)

//单一post
type Post struct {
	Tid       int //id
	Cid       int
	Uid       int
	Username  string
	Title     string
	Content   string
	Tags      string
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
	Post     *Post //id
	Comments []*Comment
}

//发布主题
func AddPost(cid, uid int, title, content string) (int64, error) {
	res, err := db.Exec(
		"call post_add(?,?,?,?)", cid, uid, title, content)

	if err != nil {
		return 0, err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return 0, code.ErrNoInsert
	}

	return rowCnt, err
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
		return code.ErrNoInsert
	}
	return err
}

//编辑文章
func EditPost(tid int, title, content string) error {
	res, err := db.Exec("call post_edit(?,?,?)",
		tid, title, content)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	log.Println("aff", rowCnt)
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return code.ErrNoInsert
	}
	return err

}

//禁止回复文章
func PostCloseComment(tid int) error {
	res, err := db.Exec("call post_close_c(?)",
		tid)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	log.Println("aff", rowCnt)
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return code.ErrNoInsert
	}
	return err
}

//允许回复文章
func PostOpenComment(tid int) error {
	res, err := db.Exec("call post_open_c(?)",
		tid)
	if err != nil {
		return err
	}
	rowCnt, err := res.RowsAffected()
	log.Println("aff", rowCnt)
	if err != nil {
		return err
	} else if rowCnt < 1 {
		return code.ErrNoInsert
	}
	return err
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

//获得文章带回复
func GetArticle(tid int) (*Article, error) {
	row := db.QueryRow(
		"SELECT `cid`,`uid`,`author`,`title`, `content`," +
			"`type`,`status`,`created`, `updated`," +
			"`views`, `replys`,`lastreply`" +
			" FROM `post` WHERE `tid` = ?", tid)

	post := &Post{Tid: tid}
	err := row.Scan(
		&post.Cid, &post.Uid, &post.Username, &post.Title, &post.Content,
		&post.Type, &post.Status, &post.Created, &post.Updated,
		&post.Views, &post.Replys, &post.Lastreply)

	//sql.ErrNoRows
	if err != nil {
		return nil, err
	}

	comments, err := GetComments(tid)
	if err != nil {
		log.Fatal(err)
		comments = nil
	}

	artcle := &Article{Post: post, Comments: comments}
	return artcle, nil
}

//按照创建顺序倒叙
func GetPostNewCreate(cid, limit, offset int) ([]*Post, error) {
	return GetPosts(cid, limit, offset, "create")
}

//按照最后回复顺序倒叙
func GetPostNewReply(cid, limit, offset int) ([]*Post, error) {
	return GetPosts(cid, limit, offset, "new")
}

//按照最近7天热帖排序
func GetPostHot(cid, limit, offset int) ([]*Post, error) {
	return GetPosts(cid, limit, offset, "hot")
}

//获得某一cid的文章列表 按照发布时间倒叙
func GetPosts(cid, limit, offset int, order string) ([]*Post, error) {
	//查询数据
	where := "ORDER BY `tid` DESC"
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
		"SELECT `tid`,`uid`,`author`,`title`, `content`," +
			"`type`,`status`,`created`, `updated`," +
			"`views`, `replys`,`lastreply`" +
			" FROM `post` WHERE `cid` = ? " + where + " LIMIT ? OFFSET ?",
		cid, limit, offset)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	posts := make([]*Post, limit)

	for rows.Next() {
		post := &Post{Cid: cid}
		err = rows.Scan(
			&post.Tid, &post.Uid, &post.Username, &post.Title, &post.Content,
			&post.Type, &post.Status, &post.Created, &post.Updated,
			&post.Views, &post.Replys, &post.Lastreply)

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
