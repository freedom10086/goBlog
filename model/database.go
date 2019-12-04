package model

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"goBlog/logger"
	"log"
)

var (
	db *sql.DB
)

//todo 要不要缓存一些 *sql.Stmt 避免不断的创建销毁

func InitDB(dbHost string, dbPort int, dbname, dbuser, dbpass string) {
	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbuser, dbpass, dbHost, dbPort, dbname)
	logger.I("connecting to %s", url)
	d, err := sql.Open("postgres", url)
	if err != nil {
		logger.E("error open db %s", url)
		log.Fatal(err)
	}

	if err = d.Ping(); err != nil {
		logger.E("error connect to db %s", url)
		log.Fatal(err)
	}

	if d != nil {
		db = d
		logger.I("success connected to the database %s", url)
	}
}

func add(sql string, args ...interface{}) (int, error) {
	var id int
	if err := db.QueryRow(sql, args...).Scan(&id); err != nil {
		return -1, err
	} else {
		return id, nil
	}
}

func del(sql string, args ...interface{}) (int64, error) {
	if res, err := db.Exec(sql, args...); err != nil {
		return -1, err
	} else {
		return res.RowsAffected()
	}
}

func update(sql string, args ...interface{}) (int64, error) {
	/*var stmt *sql.Stmt
	if stmt, err = db.Prepare(sql); err != nil {
		return
	}
	defer stmt.Close()
	var res sql.Result
	if res, err = stmt.Exec(args...); err != nil {
		return
	}
	*/
	if res, err := db.Exec(sql, args); err != nil {
		return -1, err
	} else {
		return res.RowsAffected()
	}
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("db close")
	}
}
