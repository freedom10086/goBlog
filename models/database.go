package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"goBlog/code"
)

var (
	db             *sql.DB
)

func InitDB(dbname, dbuser, dbpass string) {
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

func delete(stmt *sql.Stmt, args ...interface{}) (deletes int64, err error) {
	deletes = -1
	res, err := stmt.Exec(args...)
	if err != nil {
		return
	}

	deletes, err = res.RowsAffected()
	if err != nil {
		return
	}

	return
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
	log.Println(affect)

	if affect <= 0 {
		return code.ErrNoUpdate
	}
	return nil
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("db close")
	}
}
