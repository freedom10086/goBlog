package models

import (
	"database/sql"
	"errors"
	"goweb/conf"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var ErrNoAff = errors.New("Insert error no data affected!")

func InitDB() {
	var err error
	//[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	db, err = sql.Open("mysql", conf.DbUsername+":"+conf.DbPassword+"@/"+conf.DbName+"?charset=utf8mb4&parseTime=true")
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
