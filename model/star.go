package model

import (
	"log"
	"time"
)

type Star struct {
	Id      int       //id
	Uid     int       //用户id
	Tid     int       //收藏帖子id
	Title   string    //收藏帖子标题
	Created time.Time //时间
}

//收藏文章
func AddStar(uid, tid int) error {
	res, err := db.Exec(
		"call star_add(?,?)", uid, tid)

	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return code.ErrNoInsert
	}
	return err
}

//取消收藏文章
func DelStarById(id int) error {
	res, err := db.Exec(
		"call star_del_byid(?)", id)

	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return code.ErrNoInsert
	}
	return err
}

//取消收藏文章
func DelStarByTid(uid, id int) error {
	res, err := db.Exec(
		"call star_del_bytid(?,?)", uid, id)

	if err != nil {
		return err
	}

	rowCnt, err := res.RowsAffected()
	if err != nil && rowCnt < 1 {
		return code.ErrNoInsert
	}
	return err
}

//获得收藏列表
func GetStars(uid int) ([]*Star, error) {
	rows, err := db.Query(
		"SELECT `id`,`tid`,`title`, `created` FROM `star` WHERE `uid` = ? ORDER BY `id` DESC", uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	starts := make([]*Star, 0)

	for rows.Next() {
		star := &Star{Uid: uid}
		err = rows.Scan(&star.Id, &star.Tid, &star.Title, &star.Created)

		if err != nil {
			log.Fatal(err)
			continue
		}
		starts = append(starts, star)
	}

	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}
	return starts, err
}
