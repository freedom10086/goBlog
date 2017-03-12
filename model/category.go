package model

import (
	"log"
	"time"
)

type Category struct {
	Cid         int
	Name        string
	Description string
	Posts       int
	Sticks      []int64
	ToadyPosts  int
	Created     time.Time
}

//新增category
func AddCate(name, description string) (int64, error) {
	s:="INSERT INTO `cate` (`name`, `description`) VALUES (?, ?)"
	return add(s, name, description)
}

//删除category 里面的帖子怎么办？
//所以最好不要删除category
func DelCate(cid int) (int64, error) {
	return del("DELETE FROM `cate` WHERE `cid` = ?", cid)
}

//获得category
func GetCate(cid int) (*Category, error) {
	cate := &Category{Cid: cid}

	err := db.QueryRow(
		"SELECT  `title`, `description`,`posts`,`todayposts`,`lastpost`,`created` FROM `category` WHERE `cid` = ?",
		cid).Scan(&cate.Name, &cate.Description, &cate.Posts, &cate.ToadyPosts, &cate.Created)
	if err != nil {
		return nil, err
	}
	return cate, err
}

//获得所有category
func GetCates() ([]*Category, error) {
	//查询数据
	rows, err := db.Query("SELECT `cid`, `name`, `description`,`posts`,`todayposts`,`created` FROM `cate`")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cates := make([]*Category, 0)

	for rows.Next() {
		cate := &Category{}
		err = rows.Scan(&cate.Cid, &cate.Name, &cate.Description, &cate.Posts, &cate.ToadyPosts, &cate.Created)
		if err != nil {
			log.Println(err)
			continue
		}
		cates = append(cates, cate)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
	}
	return cates, err
}

//修改category
func ModifyCate(cid int, name, description string) (int64, error) {
	sql := "UPDATE `category` SET `title` = ?,`description` = ? WHERE `cid` = ?"
	return update(sql, name, description, cid)
}
