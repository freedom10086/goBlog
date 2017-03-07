package models

import (
	"log"
	"time"
	"goBlog/code"
)

type Comment struct {
	Id      int
	User    User
	Tid     int
	Pid     int
	Tuid    int
	Content string
	Created time.Time
	Updated time.Time
	IsRead  int
	Replys  int
}

//发表回复回复楼主
func AddCommentLz(tid, uid int, content string) error {
	if PostCanReply(tid) {
		res, err := db.Exec("call comment_add_lz(?,?,?)",
			tid, uid, content)
		if err != nil {
			return err
		}

		rowCnt, err := res.RowsAffected()
		if err != nil && rowCnt < 1 {
			return code.ErrNoInsert
		}
		return err
	} else {
		return code.ErrReply
	}
}

//发表回复回复层主
func AddCommentCz(tid, pid, uid int, content string) error {
	if PostCanReply(tid) {
		res, err := db.Exec("call comment_add_cz(?,?,?,?)",
			tid, pid, uid, content)
		if err != nil {
			return err
		}

		rowCnt, err := res.RowsAffected()
		if err != nil && rowCnt < 1 {
			return code.ErrNoInsert
		}
		return err
	} else {
		return code.ErrReply
	}
}

//获得评论
func GetComment(id int) (*Comment, error) {
	comment := &Comment{Id: id}
	err := db.QueryRow(
		"SELECT  `tid`,`pid`,`uid`,`tuid`,`author`,`content`,`created`,`updated`,`isread`,`replys` FROM `comment` WHERE `id` = ?",
		id).Scan(&comment.Tid, &comment.Pid, &comment.User.Uid, &comment.Tuid, &comment.User.Username,
		&comment.Content, &comment.Created, &comment.Updated, &comment.IsRead, &comment.Replys)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

//获得某一文章的所有评论/或者某一评论的子评论
func GetComments(tid int) ([]*Comment, error) {

	rows, err := db.Query(
		"SELECT `id`,`pid`,`uid`,`tuid`,`author`,`content`,`created`,`updated`,`isread`,`replys` FROM `comment` WHERE `tid` = ?", tid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := make([]*Comment, 10)
	for rows.Next() {
		comment := &Comment{Tid: tid}
		err = rows.Scan(&comment.Id, &comment.Pid, &comment.User.Uid, &comment.Tuid, &comment.User.Username,
			&comment.Content, &comment.Created, &comment.Updated, &comment.IsRead, &comment.Replys)
		if err != nil {
			log.Fatal(err)
			continue
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return comments, err
}

//获得某一楼楼中楼评论
func GetCommentsLzl(tid, id int) ([]*Comment, error) {
	//todo	是不是不用tid查询更快？？
	rows, err := db.Query(
		"SELECT `id`,`uid`,`tuid`,`author`,`content`,`created`,`updated`,`isread`,`replys` FROM `comment` WHERE `tid` = ? AND `pid` = ?", tid, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := make([]*Comment, 3)
	for rows.Next() {
		comment := &Comment{Tid: tid, Pid: id}
		err = rows.Scan(&comment.Id, &comment.User.Uid, &comment.Tuid, &comment.User.Username,
			&comment.Content, &comment.Created, &comment.Updated, &comment.IsRead, &comment.Replys)
		if err != nil {
			log.Fatal(err)
			continue
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return comments, err
}

//删除回复
func DelComment(id int) error {
	res, err := db.Exec(
		"call comment_del(?);", id)
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

//修改回复
func ModifyComment(id int, content string) error {

	res, err := db.Exec("call comment_edit(?,?)", id, content)

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

//回复消息已读
func SetReadComment_s(id int) error {
	res, err := db.Exec("call comment_read_s(?)", id)

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

//某篇回复消息全部已读
func SetReadComment_t(uid, tid int) error {
	res, err := db.Exec("call comment_read_t(?,?)", uid, tid)

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

//某人回复消息全部已读
func SetReadComment_a(uid int) error {
	res, err := db.Exec("call comment_read_a(?)", uid)

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
