package main

import "time"

type Book struct {
	Id          int64     `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	PublishedAt time.Time `json:"published_at"`
}

type BookInput struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishedAt string `json:"published_at"`
}

type JsonResponse map[string]any
