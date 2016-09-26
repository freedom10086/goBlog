package models

import (
	"log"
)

/*

DROP TABLE IF EXISTS `category`

CREATE TABLE IF NOT EXISTS `category` (
        `cid` integer AUTO_INCREMENT NOT NULL PRIMARY KEY,
        `name` varchar(50) NOT NULL ,
        `description` varchar(100) NOT NULL DEFAULT '' ,
        `posts` integer NOT NULL DEFAULT 0
    ) ENGINE=InnoDB;
    CREATE INDEX `category_name` ON `category` (`name`);
*/

type Category struct {
	Cid         int
	Name        string
	Description string
	Posts       int
}

//新增category
func AddCate(name, description string) error {
	_, err := db.Exec(
		"INSERT INTO `category` (`name`,`description`) VALUES (?,?)",
		name, description)
	return err
}

//删除category 里面的帖子怎么办？
//所以最好不要删除category
func DelCate(cid int) error {
	_, err := db.Exec(
		"DELETE FROM `category` WHERE cid = ?",
		cid)
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM `post` WHERE `cid` = ?")
	return err
}

//获得category
func GetCate(cid int) (*Category, error) {
	//查询数据
	row := db.QueryRow("SELECT  `name`, `description` FROM `category` WHERE `cid` = ?",
		cid)

	cate := &Category{Cid: cid}

	err := row.Scan(&cate.Name, &cate.Description, &cate.Posts)

	if err == nil {
		return cate, nil
	}

	return nil, err
}

//获得所有category
func GetCates() ([]*Category, error) {
	//查询数据
	rows, err := db.Query("SELECT `cid`, `name`, `description`, `posts` FROM `category` ", nil)

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	cates := make([]*Category, 0)

	for rows.Next() {
		cate := &Category{}
		err = rows.Scan(&cate.Cid, &cate.Name, &cate.Description, &cate.Posts)
		if err != nil {
			log.Fatal(err)
			continue
		}
		cates = append(cates, cate)
	}

	return cates, err
}

//修改category
func ModifyCate(cid int, name, description string) error {

	_, err := db.Exec(
		"UPDATE `category` SET  `name` = ?, `description` = ? WHERE `tid` = ?",
		name, description, cid)

	return err
}
