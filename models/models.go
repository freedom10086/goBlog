package models

import (
	"database/sql"
	"errors"
	"log"

	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var (
	db          *sql.DB
	ErrNoInsert error = errors.New("操作失败!")
	ErrNoDelete error = errors.New("删除失败!")
	ErrNoUpdate error = errors.New("没有更改!")
	ErrReply    error = errors.New("此文章无法回复!")
	ErrLogin          = errors.New("账号异常,你没有权限登陆!")
)

func InitDB(dbuser, dbpass, dbname string) {
	var err error
	db, err = sql.Open("mysql", dbuser+":"+dbpass+"@/"+dbname+"?charset=utf8mb4&parseTime=true")
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	if db != nil {
		log.Println("sql connect success")
	}
}

func add(stmt *sql.Stmt, args ...interface{}) (id int64, err error) {
	id = -1
	res, err := stmt.Exec(args...)
	if err != nil {
		return
	}

	id, err = res.LastInsertId()
	if err != nil {
		return
	}
	return
}

func delete(stmt *sql.Stmt, args ...interface{}) error {
	res, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(affect)

	if affect <= 0 {
		return ErrNoDelete
	}
	return nil
}

func update(stmt *sql.Stmt, args ...interface{}) error {
	res, err := stmt.Exec(args...)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}
	fmt.Println(affect)

	if affect <= 0 {
		return ErrNoUpdate
	}
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("db close")
	}
}
