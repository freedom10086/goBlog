package model

import (
	"time"
)

type Chat struct {
	Id       int
	Uid      int
	Tuid     int
	Username string `json:",omitempty"` //N 一直是对方用户名，不是tuid的用户名，有可能tuid是自己
	Content  string
	IsRead   int       //是否已读 0 未读
	Created  time.Time //时间
}

//发送私信
func AddChat(uid, tuid int, content string) (int64, error) {
	if uid == tuid {
		return -1, ErrParama
	}
	sql := "INSERT INTO `chat` (`uid`, `tuid`,`content`) VALUES (?,?,?)"
	return add(sql, uid, tuid, content)
}

//撤回
//1分钟内没读可以撤回
func DelChat(id int) (int64, error) {
	sql := "delete from `chat` where id = ? and timestampdiff(minute,now(),created) <1 and isread = 0"
	return del(sql, id)
}

//获得我和另一个人的对话
func GetChats(uid, tuid, page, pageSize int) (cs []*Chat, err error) {
	offset := (page - 1) * pageSize
	s := `
	SELECT id,uid,tuid,content,created FROM chat
	WHERE
	(uid = ? AND tuid = ?)
	OR
	(tuid = ? AND uid = ?)
	ORDER BY id DESC LIMIT ? OFFSET ?`
	rows, err := db.Query(s, uid, tuid, uid, tuid, pageSize, offset)
	if err != nil {
		return
	}
	defer rows.Close()
	cs = make([]*Chat, 0, pageSize)
	for rows.Next() {
		c := &Chat{}
		if err = rows.Scan(&c.Id, &c.Uid, &c.Tuid, &c.Content, &c.Created); err != nil {
			return
		}
		cs = append(cs, c)
	}
	err = rows.Err()
	return
}

//获得和我聊天的列表
//要有最后一次聊天的内容放在列表
//类似于qq消息列表
//sender 1-表示这条消息是发送方是我
//       0-表示这条消息接收方是我
func GetRecentChats(uid, page, pageSize int) (cs []*Chat, err error) {
	s := `
	SELECT t.id,t.ouid,u.username,t.content,t.sender,t.created
	FROM (
	SELECT id,ouid,content,sender,created FROM
		((SELECT id,tuid as ouid,content, 1 as sender,created FROM chat WHERE uid = ?)
	 	UNION
	 	(SELECT id,uid as ouid,content,0 as sender,created FROM chat WHERE tuid = ?)
	 	ORDER BY id DESC)
		as tmp
		GROUP BY tmp.ouid
		ORDER BY tmp.id DESC
		limit ? OFFSET ?
	) as t
	LEFT JOIN user AS u ON t.ouid = u.uid`
	offset := (page - 1) * pageSize
	rows, err := db.Query(s, uid, uid, pageSize, offset)
	if err != nil {
		return
	}
	defer rows.Close()
	cs = make([]*Chat, 0, pageSize)
	var sender int
	for rows.Next() {
		c := &Chat{Uid: uid}
		if err = rows.Scan(&c.Id, &c.Tuid, &c.Username, &c.Content, &sender, &c.Created); err != nil {
			return
		}
		//tuid 是我
		if sender == 0 {
			c.Uid = c.Tuid
			c.Tuid = uid
		}
		cs = append(cs, c)
	}
	err = rows.Err()
	return

}
