package main

import (
	"embed"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"time"
)

// Шаблоны и CSS включаются в программу при сборке.
//
//go:embed templates/*.html static/*.css
var assets embed.FS

var pages = template.Must(template.ParseFS(assets, "templates/*.html"))

func routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/tasks", addTaskHandler)
	mux.HandleFunc("/tasks/complete", completeTaskHandler)
	mux.HandleFunc("/tasks/delete", deleteTaskHandler)
	mux.HandleFunc("/about", aboutHandler)
	staticFiles, err := fs.Sub(assets, "static")
	if err != nil {
		panic(err)
	}
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFiles))))
	return mux
}

func main() {
	server := &http.Server{
		Addr:              "localhost:8080",
		Handler:           routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Println("Task Tracker: http://localhost:8080 — откройте адрес в браузере. Остановка: Ctrl+C.")
	log.Fatal(server.ListenAndServe())
}
