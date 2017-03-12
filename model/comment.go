package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	Id      int
	Tid     int
	Pid     int
	Uid     int
	Tuid    int
	Content string
	Replys  int
	Created time.Time
	Updated time.Time
}

//发表回复回复楼主
func AddCommentLz(tid, uid int, content string) (int64, error) {
	s := "call comment_add_lz(?,?,?)"
	return add(s, tid, uid, content)
}

//发表回复回复层主
func AddCommentCz(tid, pid, uid int, content string) (int64, error) {
	if status, err := getPostStatus(tid); err != nil {
		return -1, err
	} else if status != 0 {
		return -1, ErrNoAuth
	}

	s := "call comment_add_cz(?,?,?,?)"
	return add(s, tid, pid, uid, content)
}

//删除回复
func DelComment(id int) (int64, error) {
	s := "call comment_del(?)"
	return del(s, id)
}

//修改回复
func UpdateComment(id int, content string) (int64, error) {
	s := "UPDATE `comment` SET `content` = ?, `updated` = now() WHERE `id` = ?"
	return update(s, content, id)
}

//获得评论
func GetComment(id int) (*Comment, error) {
	c := &Comment{Id: id}
	s := "SELECT  `tid`,`pid`,`uid`,`tuid`,`content`,`replys`,`created`,`updated` FROM `comment` WHERE `id` = ?"
	err := queryA1(s, id, &c.Tid, &c.Pid, &c.Uid, &c.Tuid, &c.Content, &c.Replys, &c.Created, &c.Updated)
	return c, err
}

//获得某一文章的所有评论
func GetComments(tid int, page, pagesize int) (cs []*Comment, err error) {
	s := "SELECT `id`,`pid`,`uid`,`tuid`,`content`,`replys`,`created`,`updated` " +
		"FROM `comment` WHERE `tid` = ? ORDER BY tid ASC LIMIT ? OFFSET ?"
	offset := (page - 1) * pagesize
	var rows *sql.Rows
	if rows, err = db.Query(s, tid, pagesize, offset); err != nil {
		return
	}
	defer rows.Close()

	cs = make([]*Comment, 0, pagesize)
	for rows.Next() {
		c := &Comment{Tid: tid}
		err = rows.Scan(&c.Id, &c.Pid, &c.Uid, &c.Tuid, &c.Content,
			&c.Replys, &c.Created, &c.Updated)
		if err != nil {
			return
		}
		cs = append(cs, c)
	}
	err = rows.Err()
	return
}

//获得某一楼楼中楼评论
func GetCommentsLzl(tid, id int) (css []*Comment, err error) {
	var rows *sql.Rows
	if rows, err = db.Query(
		"SELECT  `id`,`uid`,`tuid`,`content`,`replys`,`created`,`updated` "+
			"FROM `comment` WHERE  `tid` = ? AND `pid` = ?", tid, id); err != nil {
		return
	}
	defer rows.Close()

	css = make([]*Comment, 0, 3)
	for rows.Next() {
		c := &Comment{Tid: tid, Pid: id}
		err = rows.Scan(&c.Id, &c.Uid, &c.Tuid, &c.Content,
			&c.Replys, &c.Created, &c.Updated)
		if err != nil {
			return
		}
		css = append(css, c)
	}
	err = rows.Err()
	return
}