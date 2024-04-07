package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	PORT = "3000"

	DB_USER     = "db_username"
	DB_PASS     = "db_password"
	DB_DATABASE = "bookstore"
)

var Db *sql.DB

func main() {
	initializeDb()

	r := mux.NewRouter()

	r.HandleFunc("/books", GetBooks).Methods(http.MethodGet)
	r.HandleFunc("/books/{id}", GetBook).Methods(http.MethodGet)
	r.HandleFunc("/books", CreateBook).Methods(http.MethodPost)
	r.HandleFunc("/books/{id}", UpdateBook).Methods(http.MethodPatch, http.MethodPut)
	r.HandleFunc("/books/{id}", DeleteBook).Methods(http.MethodDelete)

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			next.ServeHTTP(w, r)
		})
	})

	fmt.Printf("Server listening on %q\n", PORT)
	fmt.Printf("Visit the URL http://127.0.0.1:%s to check if the connection was successful\n", PORT)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", PORT), r))
}

func initializeDb() {
	dsn := fmt.Sprintf("%s:%s@/%s?parseTime=true", DB_USER, DB_PASS, DB_DATABASE)

	var err error
	Db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalln("Error connecting to database -", err.Error())
	}

	if err = Db.Ping(); err != nil {
		log.Fatalln("Error pinging database -", err.Error())
	}

	stmt := `
	CREATE TABLE IF NOT EXISTS books (
		id INT AUTO_INCREMENT PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(255) NOT NULL,
		published_at DATETIME NOT NULL
	);`

	if _, err := Db.Query(stmt); err != nil {
		log.Fatalln("Error creating table -", err.Error())
	}
}
