package model

import "time"

type Follow struct {
	Id       int
	Uid      int
	Tuid     int
	Username string `json:",omitempty"` //对方的用户名 //可能是uid的也可能是tuid的
	Note     string `json:",omitempty"`
	Created  time.Time
}

//关注用户
func AddFollow(uid, tuid int) (int64, error) {
	if uid == tuid {
		return -1, ErrParama
	}
	sql := "INSERT INTO `follow` (`uid`, `tuid`) VALUES ($1,$2)"
	return add(sql, uid, tuid)
}

//取消follow
func DelFollowById(id int) (int64, error) {
	sql := "delete from `follow` where id = $1"
	return del(sql, id)
}

//取消关注
func DelStarByUid(uid, tuid int) (int64, error) {
	sql := "delete from `follow` where uid = $1 and tuid = $2"
	return del(sql, uid, tuid)
}

//获得我关注列表
func GetFollows(uid, page, pageSize int) (follows []*Follow, err error) {
	offset := (page - 1) * pageSize
	s := `
	SELECT f.id,f.tuid,u.username,f.note,f.created
	FROM follow as f,user as u
	WHERE f.uid = $1 AND f.tuid = u.uid
	ORDER BY f.id desc
	LIMIT $2 OFFSET $3`

	rows, err := db.Query(s, uid, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	follows = make([]*Follow, 0, pageSize)
	for rows.Next() {
		f := &Follow{Uid: uid}
		if err = rows.Scan(&f.Id, &f.Tuid, &f.Username, &f.Note, &f.Created); err != nil {
			return
		}
		follows = append(follows, f)
	}

	err = rows.Err()
	return
}

//获得关注我的
func GetFollowsMe(uid, page, pageSize int) (follows []*Follow, err error) {
	offset := (page - 1) * pageSize
	s := `
	SELECT f.id,f.uid,u.username,f.created
	FROM follow as f,user as u
	WHERE f.tuid = $1 AND f.uid = u.uid
	ORDER BY f.id desc
	LIMIT $2 OFFSET $3`
	rows, err := db.Query(s, uid, pageSize, offset)
	if err != nil {
		return
	}
	defer rows.Close()
	follows = make([]*Follow, 0, pageSize)
	for rows.Next() {
		f := &Follow{Tuid: uid}
		if err = rows.Scan(&f.Id, &f.Uid, &f.Username, &f.Created); err != nil {
			return
		}
		follows = append(follows, f)
	}
	err = rows.Err()
	return
}
