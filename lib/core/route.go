package server

import "net/http"

type RouteHandler func(request *http.Request, controller *Controller) error

type Route struct {
	Method  string
	Handler RouteHandler
	Path    string
}

func NewRoute(path string, handler RouteHandler) *Route {
	return &Route{
		Handler: handler,
		Path:    path,
	}
}
