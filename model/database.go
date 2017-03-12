package model

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	db          *sql.DB
	ErrNoRecord error = errors.New("没有查询到记录")
	ErrNoAffect error = errors.New("没有改变")
	ErrNoAuth error = errors.New("没有权限")
	ErrParama error = errors.New("参数错误")
)

//todo 要不要缓存一些 *sql.Stmt 避免不断的创建销毁

func InitDB(dbname, dbuser, dbpass string) {
	var err error
	db, err = sql.Open("mysql", dbuser + ":" + dbpass + "@/" + dbname + "?charset=utf8mb4&parseTime=true")
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

func add(sql string, args ...interface{}) (int64, error) {
	if res, err := db.Exec(sql, args...); err != nil {
		return -1, err
	} else {
		return res.LastInsertId()
	}
}

func del(sql string, args ...interface{}) (deletes int64, err error) {
	if res, err := db.Exec(sql, args...); err != nil {
		return -1, err
	} else {
		return res.RowsAffected()
	}
}

func update(sql string, args ...interface{}) (affected int64, err error) {
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

//一个参数的query
func queryA1(sql string, arg interface{}, dest ...interface{}) (error) {
	return db.QueryRow(sql, arg).Scan(dest)
}

//2个参数的query
func queryA2(sql string, arg1 interface{}, arg2 interface{}, dest ...interface{}) (err error) {
	return db.QueryRow(sql, arg1, arg2).Scan(dest)
}

func CloseDB() {
	if db != nil {
		db.Close()
		log.Println("db close")
	}
}
