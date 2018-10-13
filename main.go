package main

import (
	"belial/server"
)

func main() {
	srv := server.NewWebServer(":8080", 1000, 1000)
	srv.Run()
}
