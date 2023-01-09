package main

import (
	server "exporter-dev/http-server/lib/core"
	"net/http"
)

func main() {
	var instance server.ServerInterface = &server.Server{
		Port: 3000,
	}

	instance.Get(*server.NewRoute("/index", func(request *http.Request, controller *server.Controller) error {
		controller.Status(201)
		controller.Header.Add("Test", "test1")
		controller.Send("OK")
		return nil
	}))

	instance.Init()
}
