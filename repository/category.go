package repository

import (
	"log"
	"time"
)

type Category struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Posts       int       `json:"posts"`
	Sticks      []int64   `json:"-"`
	Created     time.Time `json:"created"`
}

//新增category
func AddCate(name, description string) (int, error) {
	s := "INSERT INTO cate (name, description) VALUES ($1, $2) RETURNING id"
	return add(s, name, description)
}

//删除category 里面的帖子怎么办？
//所以最好不要删除category
func DelCate(id int) (int64, error) {
	return del("DELETE FROM cate WHERE id = $1", id)
}

//获得category
func GetCate(id int) (*Category, error) {
	cate := &Category{Id: id}
	s := "SELECT name, description, posts ,created FROM category WHERE id = $1"
	err := db.QueryRow(s, id).Scan(&cate.Name, &cate.Description, &cate.Posts, &cate.Created)
	if err != nil {
		return nil, err
	}
	return cate, err
}

//获得所有category
func GetCates() ([]*Category, error) {
	//查询数据
	rows, err := db.Query("SELECT id, name, description, posts, created FROM category")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cates := make([]*Category, 0)

	for rows.Next() {
		cate := &Category{}
		err = rows.Scan(&cate.Id, &cate.Name, &cate.Description, &cate.Posts, &cate.Created)
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
func ModifyCate(id int, name, description string) (int64, error) {
	sql := "UPDATE category SET title = $1,description = $2 WHERE id = $3"
	return update(sql, name, description, id)
}
