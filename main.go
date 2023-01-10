package main

import (
	server "exporter-dev/http-server/lib/core"
	"log"
)

func main() {
	var instance server.ServerInterface = &server.Server{
		Port: 3000,
	}

	instance.Get(
		*server.NewRoute(
			"/index",
			func(request *server.Request, controller *server.Controller) error {
				controller.Status(201)
				controller.Header.Add("Test", "test1")
				controller.Send("Your IP: " + request.RemoteAddt)
				return nil
			},
		),
		*server.NewRoute(
			"/index/:param/",
			func(request *server.Request, controller *server.Controller) error {
				controller.Status(200)

				params := ""

				for key, value := range request.Params {
					params += key + ":" + value + ";"
				}

				// request.Params: map[param:1] -> URL: /index/1/
				log.Println(request.Params)

				controller.Send("Your params: " + params)
				return nil
			},
			server.RouteOptions{
				HasParams: true,
			},
		),
	)

	instance.Init()
}
