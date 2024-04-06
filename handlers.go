package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest/date"
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

func CreateBook(w http.ResponseWriter, r *http.Request) {
	var (
		book  Book
		input BookInput
	)

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "Unprocessable entity"})
		return
	}

	if book.Author = strings.TrimSpace(input.Author); len(book.Author) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "author field is required"})
		return
	}

	if book.Title = strings.TrimSpace(input.Title); len(book.Title) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "title field is required"})
		return
	}

	book.PublishedAt, err = formatDate(input.PublishedAt)
	switch {
	case err != nil:
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "error parsing published_at field - ensure it is of the format 'YYYY-MM-DD'"})
		return
	case book.PublishedAt.IsZero():
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "published_at field is not a valid date"})
		return
	}

	stmt := "INSERT INTO books (author, title, published_at) VALUES (?, ?, ?)"
	result, err := Db.Exec(stmt, book.Author, book.Title, book.PublishedAt)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{"message": "error saving book"})
		return
	}

	book.Id, _ = result.LastInsertId()

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JsonResponse{"data": book})
}

func UpdateBook(w http.ResponseWriter, r *http.Request) {

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{"message": "invalid book ID"})
		return
	}

	var (
		book  Book
		input BookInput
	)

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "Unprocessable entity"})
		return
	}

	if book.Author = strings.TrimSpace(input.Author); len(book.Author) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "author field is required"})
		return
	}

	if book.Title = strings.TrimSpace(input.Title); len(book.Title) < 1 {
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "title field is required"})
		return
	}

	book.PublishedAt, err = formatDate(input.PublishedAt)
	switch {
	case err != nil:
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "error parsing published_at field - ensure it is of the format 'YYYY-MM-DD'"})
		return
	case book.PublishedAt.IsZero():
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(JsonResponse{"message": "published_at field is not a valid date"})
		return
	}

	stmt := "UPDATE books SET author = ?, title = ?, published_at = ? WHERE id = ?"
	_, err = Db.Exec(stmt, book.Author, book.Title, book.PublishedAt, id)
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{"message": "error updating book"})
		return
	}

	book.Id = int64(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(JsonResponse{"data": book})
}

func DeleteBook(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	if id == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(JsonResponse{"message": "invalid book ID"})
		return
	}

	result, err := Db.Exec("DELETE FROM books where id = ?", id)
	rows, _ := result.RowsAffected()
	if err != nil || rows < 1 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(JsonResponse{"message": "error deleting book"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(JsonResponse{"message": "book deleted successfully"})
}

func formatDate(s string) (time.Time, error) {
	date, err := date.ParseDate(s)
	if err != nil {
		return time.Time{}, err
	}

	return date.ToTime(), nil
}
