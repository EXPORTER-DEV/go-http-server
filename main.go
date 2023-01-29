package main

import (
	server "exporter-dev/http-server/lib/core"
	"fmt"
	"log"
)

type Test struct {
	a string
	b int
}

func main() {
	var instance = &server.Server{
		Port: 3000,
	}

	instance.Use(
		*server.NewMiddleware(
			func(request *server.Request, controller *server.Controller) (skip bool, err error) {
				request.Context.Set("Test", Test{
					a: "test",
					b: 0,
				})

				if request.Headers["Auth"] != nil {
					log.Printf("User authorized!")
					request.Context.Set("Auth", true)
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

	indexRouteGroup := server.NewRouteGroup(
		"/index",
		*server.NewRoute(
			"",
			func(request *server.Request, controller *server.Controller) error {
				controller.Status(201)
				controller.Header.Add("Test", "test1")
				controller.Send("Your IP: " + request.RemoteAddt)
				return nil
			},
		),
		*server.NewRoute(
			"/:param/",
			func(request *server.Request, controller *server.Controller) error {
				test := server.BindContext[Test](request.Context, "Test")

				fmt.Printf("Got Test structure from Context: %v\n", test)

				auth := server.BindContext[bool](request.Context, "Auth")

				if !auth {
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

	instance.Get(
		indexRouteGroup...,
	)

	instance.Init()
}
