package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBooks(w http.ResponseWriter, r *http.Request) {
	rows, err := Db.Query("SELECT * FROM books")
	if err != nil {
		log.Println("Error fetching result from db -", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{"message": "Oops, something went wrong"})
		return
	}

	books := []Book{}

	for rows.Next() {
		var book Book

		err := rows.Scan(&book.Id, &book.Title, &book.Author, &book.PublishedAt)
		if err != nil {
			continue
		}

		books = append(books, book)
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{"data": books})
}

func GetBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{"message": "Invalid book ID"})
		return
	}

	var book Book

	row := Db.QueryRow("SELECT * FROM books where id = ?", id)
	err := row.Scan(&book.Id, &book.Title, &book.Author, &book.PublishedAt)

	switch {
	case err == sql.ErrNoRows:
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(JsonResponse{"message": "book not found"})
	case err != nil:
		log.Println("Error fetching book from db -", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{"message": "Oops, something went wrong"})
	default:
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(JsonResponse{"data": book})
	}
}
