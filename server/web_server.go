package server

import (
	"belial/tasks"
	"html/template"
	"log"
	"net/http"
	"time"
)

// WebServer - main web-listner
type WebServer struct {
	Addr     string       // Server address (or only port)
	RTimeout uint         // Read timeout in milliseconds
	WTimeout uint         // Write timeout in milliseconds
	Tasks    []tasks.Task // Tasks storage, temporary - slice
	server   http.Server
}

// Run - load templates and start HTTP-server
func (s *WebServer) Run() {

	tmpl := template.New("")

	_, err := tmpl.ParseFiles("templates/main.html", "templates/test.html")
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

	s.server = http.Server{
		Addr:         s.Addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(s.RTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(s.WTimeout) * time.Millisecond,
	}

	log.Println("Starting server at", s.Addr)
	s.server.ListenAndServe()
}
