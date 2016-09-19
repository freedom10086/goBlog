package models

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() {
	var err error
	db, err = sql.Open("mysql", "root:justice@/goweb?charset=utf8")
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

func CloseDB() {
	if db != nil {
		db.Close()
	}
}
