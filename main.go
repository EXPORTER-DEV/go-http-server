package main

import (
	server "exporter-dev/http-server/lib/core"
	"log"
)

func main() {
	var instance = &server.Server{
		Port: 3000,
	}

	instance.Use(
		*server.NewMiddleware(
			func(request *server.Request, controller *server.Controller) (skip bool, err error) {
				if request.Headers["Auth"] != nil {
					log.Printf("User authorized!")
					request.Context.Set("Auth", "1")
				}

				return false, nil
			},
		),
		*server.NewMiddleware(
			func(request *server.Request, controller *server.Controller) (skip bool, err error) {
				log.Println("/index only middleware")

				return false, nil
			},
		).AddPath("/index"),
	)

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
				if request.Context.Get("Auth") == "" {
					controller.Status(401)
				} else {
					controller.Status(200)
				}

				params := ""

				for key, value := range request.Params {
					params += key + ":" + value + ";"
				}

				// request.Params: map[param:1] -> URL: /index/1/
				log.Println(request.Params)

				controller.Send("Your params: " + params)
				return nil
			},
		).SetParseParams(true),
	)

	instance.Init()
}
