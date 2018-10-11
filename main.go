package main

import (
	"belial/tasks"
	"html/template"
	"log"
	"net/http"
	"time"
)

// Server - main web
type Server struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Tasks        []tasks.Task
}

func (s *Server) run() {

	tmpl := template.New("")

	_, err := tmpl.ParseFiles("templates/main.html")
	if err != nil {
		panic(err)
	}
	_, err = tmpl.ParseFiles("templates/test.html")
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "main.html", s)
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "test.html", s)
		if err != nil {
			panic(err)
		}
	})

	server := http.Server{
		Addr:         s.Addr,
		Handler:      mux,
		ReadTimeout:  s.ReadTimeout,
		WriteTimeout: s.WriteTimeout,
	}

	log.Println("Starting server at", s.Addr)
	server.ListenAndServe()
}

func main() {

	tasks := []tasks.Task{
		tasks.Task{1, "Test 1", true},
		tasks.Task{2, "Test X", false},
		tasks.Task{3, "Test 2", true},
	}

	srv := Server{
		Addr:         ":8080",
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		Tasks:        tasks,
	}

	srv.run()
}
