package main

import (
	"exporter-dev/http-server/server"
	"net/http"
)

func main() {
	var instance server.ServerInterface = &server.Server{
		Port: 3000,
	}

	instance.Get("/index", func(request *http.Request, response http.ResponseWriter) error {
		response.Write([]byte("Test"))
		return nil
	})

	instance.Init()
}
