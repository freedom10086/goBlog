package model

import (
	"time"
)

type Star struct {
	Id      int       //id
	Uid     int       //用户id
	Tid     int       //收藏帖子id
	Created time.Time //时间
}

//收藏文章
func AddStar(uid, tid int) (int, error) {
	sql := "INSERT INTO star (uid, tid) VALUES ($1,$2) RETURNING id"
	return add(sql, uid, tid)
}

//取消收藏文章
func DelStarById(id int) (int64, error) {
	sql := "delete from star where id = $1"
	return del(sql, id)
}

//取消收藏文章
func DelStarByTid(uid, tid int) (int64, error) {
	sql := "delete from star where uid = $1 and tid = $2"
	return del(sql, uid, tid)
}

//获得收藏列表
func GetStars(uid, page, pageSize int) (starts []*Star, err error) {
	offset := (page - 1) * pageSize
	rows, err := db.Query(
		"SELECT id,tid,created FROM star WHERE uid = $1 ORDER BY id DESC LIMIT $2 OFFSET $3",
		uid, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	starts = make([]*Star, 0, pageSize)
	for rows.Next() {
		star := &Star{Uid: uid}
		if err = rows.Scan(&star.Id, &star.Tid, &star.Created); err != nil {
			return
		}
		starts = append(starts, star)
	}

	err = rows.Err()
	return
}
