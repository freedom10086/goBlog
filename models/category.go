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

func AddCate(name, description string) error {
	_, err := db.Exec(
		"INSERT INTO `category` (`name`,`description`) VALUES (?,?)",
		name, description)
	return err
}

func DelCate(cid int) error {
	_, err := db.Exec(
		"DELETE FROM `category` WHERE cid = ?",
		cid)
	return err
}

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

func ModifyCate(cid int, name, description string) error {

	_, err := db.Exec(
		"UPDATE `category` SET  `name` = ?, `description` = ? WHERE `tid` = ?",
		name, description, cid)

	return err
}
