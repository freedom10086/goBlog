package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	var err error
	s := config.DbUsername + ":" + config.DbPassword + "@/" + config.DbName
	//root:justice@/goblog
	db, err = sql.Open("mysql", s)
	if err != nil {
		panic(err.Error())
	}

	log.Println("success connected to db!");
	defer db.Close()
}
