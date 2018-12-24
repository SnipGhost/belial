package server

import (
	"belial/tasks"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// WebServer - main web-listner
type WebServer struct {
	Addr     string         // Server address, or port only
	RTimeout uint           // Read timeout in milliseconds
	WTimeout uint           // Write timeout in milliseconds
	Data     *tasks.Storage // Tasks storage
	server   *http.Server
}

// NewWebServer - return initialized WebServer
func NewWebServer(addr string, rtimeout uint, wtimeout uint) WebServer {

	websrv := WebServer{
		Addr:     addr,
		RTimeout: rtimeout,
		WTimeout: wtimeout,
		Data:     tasks.NewStorage(),
	}

	//log.Printf("Created WebServer: %+v", websrv)
	return websrv
}

// Run - load templates and start HTTP-server
func (s *WebServer) Run() {

	tmpl := template.New("")

	_, err := tmpl.ParseFiles(
		"templates/main.html",
		"templates/test.html",
		"templates/form.html",
	)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "main.html", s.Data)
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/add/", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		data := r.FormValue("data")
		n := r.FormValue("n")
		k := r.FormValue("k")
		timeout, err := strconv.ParseUint(r.FormValue("timeout"), 10, 64)
		if err == nil && len(data) > 0 && len(n) > 0 && len(k) > 0 {
			s.Data.AddTask(tasks.NewTask(data, n, k, "0", timeout))
		}
		err = tmpl.ExecuteTemplate(w, "form.html", s.Data)
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

	s.server = &http.Server{
		Addr:         s.Addr,
		Handler:      mux,
		ReadTimeout:  time.Duration(s.RTimeout) * time.Millisecond,
		WriteTimeout: time.Duration(s.WTimeout) * time.Millisecond,
	}

	go s.Data.Scheduler.Run()

	log.Println("Starting web-server at", s.Addr, "...")
	err = s.server.ListenAndServe()
	if err != nil {
		log.Println("Could not starting web-server:", err)
	}
}
