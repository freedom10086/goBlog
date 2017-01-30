package main

import (
	"fmt"
	"net/http"
	"database/sql"
	"goweb/handlers"
	"log"
	_ "github.com/lib/pq"
)

var db *sql.DB

type myHandler struct {
}

func (*myHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	url := "https://127.0.0.1:10443" + req.RequestURI
	http.Redirect(w, req, url, http.StatusMovedPermanently)
}

func init() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=justice dbname=postgres sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	queryMulRows()

	querySingle()

	go func() {
		log.Printf("About to listen on 8080. Go to http://127.0.0.1:8080/")
		http.ListenAndServe(":8080", &myHandler{})
	}()

	http.Handle("/", handlers.NewStaticServer())
	log.Printf("About to listen on 10443. Go to https://127.0.0.1:10443/")
	err := http.ListenAndServeTLS(":10443", "cert.pem", "key.pem", nil)
	log.Fatal(err)
}

func queryMulRows() {
	fmt.Println("test multi query")
	rows, err := db.Query("SELECT username FROM public.user")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", name)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func querySingle() {
	fmt.Println("test single query")
	var uid int32 = 1
	row := db.QueryRow("SELECT username FROM public.user WHERE uid = $1", uid)
	var name string
	err := row.Scan(&name)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)
}

func queryWithNull() {
	/*
	var name sql.NullString
	err := db.QueryRow("SELECT name FROM names WHERE id = $1", id).Scan(&name)
	if name.Valid {
		// use name.String
	} else {
		// value is NULL
	}
	*/
}

func dealLastInserId() {
	/*
	var userid int
	err := db.QueryRow(`INSERT INTO users(name, favorite_fruit, age)
	VALUES('beatrice', 'starfruit', 93) RETURNING id`).Scan(&userid)
	*/

}

/*
- integer types smallint, integer, and bigint are returned as int64
- floating-point types real and double precision are returned as float64
- character types char, varchar, and text are returned as string
- temporal types date, time, timetz, timestamp, and timestamptz are returned as time.Time
- the boolean type is returned as bool
- the bytea type is returned as []byte
*/
