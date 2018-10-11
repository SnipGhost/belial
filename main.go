package main

import (
	"belial/server"
	"belial/tasks"
)

func main() {

	tasks := []tasks.Task{
		tasks.Task{1, "Test 1", true},
		tasks.Task{2, "Test X", false},
		tasks.Task{3, "Test 2", true},
		tasks.Task{4, "Test", true},
	}

	srv := server.WebServer{
		Addr:     ":8080",
		RTimeout: 1000,
		WTimeout: 1000,
		Tasks:    tasks,
	}

	srv.Run()
}
