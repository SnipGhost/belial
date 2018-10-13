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
	Addr     string         // Server address, or port only
	RTimeout uint           // Read timeout in milliseconds
	WTimeout uint           // Write timeout in milliseconds
	Data     *tasks.Storage // Tasks storage
	server   *http.Server
}

// NewWebServer - return initialized WebServer
func NewWebServer(addr string, rtimeout uint, wtimeout uint) WebServer {

	storage := tasks.NewStorage()

	/* START OF TESTING BLOCK */
	var tps []*tasks.Task
	var ids []uint

	tps = append(tps, tasks.NewTask("Test-1", 100))   // 0
	tps = append(tps, tasks.NewTask("Test-2", 10000)) // 1
	tps = append(tps, tasks.NewTask("Test-3", 10000)) // 2
	tps = append(tps, tasks.NewTask("Test-4", 20000)) // 3

	for _, tp := range tps {
		ids = append(ids, storage.AddTask(tp))
	}

	storage.CancelTask(ids[2])
	/* END OF TESTING BLOCK */

	websrv := WebServer{
		Addr:     addr,
		RTimeout: rtimeout,
		WTimeout: wtimeout,
		Data:     storage,
	}

	//log.Printf("Created WebServer: %+v", websrv)
	return websrv
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
		err := tmpl.ExecuteTemplate(w, "main.html", s.Data)
		if err != nil {
			panic(err)
		}
	})
	mux.HandleFunc("/test/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "test.html", s.Data)
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
