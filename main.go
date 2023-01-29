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
	// Create new Server instance, and SetPort (by default will use 80 port):
	var instance = server.NewServer().SetPort(3000)

	// Register middlewares in first register - first execute order:
	instance.Use(
		*server.NewMiddleware(
			func(request *server.Request, controller *server.Controller) (skip bool, err error) {
				// Set Context value by "Test" key with Test struct:
				request.Context.Set("Test", Test{
					a: "test",
					b: 0,
				})

				if request.Headers["Auth"] != nil {
					log.Printf("User authorized!")
					// Set Context value by "Auth" key with boolean value:
					request.Context.Set("Auth", true)
				}

				return false, nil
			},
		),
		// Register middleware and specify the path when it will be triggered:
		*server.NewMiddleware(
			func(request *server.Request, controller *server.Controller) (skip bool, err error) {
				log.Println("/index only middleware")

				// Each middleware handler should return (skip bool, err error),
				// when skip is true, then next execution stops:
				return false, nil
			},
		).AddPath("/index"),
	)

	// Create RouteGroup with path where all children routes
	// will be accessed with "/index" + children route path:
	indexRouteGroup := server.NewRouteGroup(
		"/index",
		*server.NewRoute(
			"",
			func(request *server.Request, controller *server.Controller) error {
				// Set controller response status to 201:
				controller.Status(201)
				// Add response header:
				controller.Header.Add("Test", "test1")
				// Send response to the client:
				controller.Send("Your IP: " + request.RemoteAddt)
				// Each route handler should return error:
				return nil
			},
		),
		// Register params Route, when use params parsing we should call SetParseParams(true):
		*server.NewRoute(
			"/:param/",
			func(request *server.Request, controller *server.Controller) error {
				// Bind context value by key and generic passed struct Test:
				test := server.BindContext[Test](request.Context, "Test")

				fmt.Printf("Got Test structure from Context: %v\n", test)

				// Bind context value by key and generic passed bool:
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

	// Register route group
	instance.Get(
		indexRouteGroup...,
	)

	// Init server listening
	instance.Init()
}
