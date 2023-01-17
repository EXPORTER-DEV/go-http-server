package server

import (
	"log"
	"net/http"
)

type RouteHandler func(request *Request, controller *Controller) error

type Route struct {
	Method      string
	Handler     RouteHandler
	ContentType string
	Path        string
	IsRegexp    bool
	ParseParams bool
}

func NewRoute(path string, handler RouteHandler) *Route {
	return &Route{
		Handler: handler,
		Path:    path,
	}
}

func NewRouteGroup(path string, routes ...Route) []Route {
	for index := range routes {
		routes[index].Path = path + routes[index].Path
	}

	return routes
}

func (r *Route) SetIsRegexp(value bool) *Route {
	r.IsRegexp = value

	return r
}

func (r *Route) SetParseParams(value bool) *Route {
	r.ParseParams = value

	return r
}

type MatchedRoute struct {
	Route  *Route
	Params Params
}

type Routing struct {
	Routes []Route
}

func (r *Routing) Match(method string, path string) *MatchedRoute {
	for i, route := range r.Routes {
		if route.Method == method {
			matching := &Matching{
				RequestedPath: path,
				HandlerPath:   route.Path,
			}

			result, params := matching.Execute(MatchingExecuteOptions{
				ParseParams: route.ParseParams,
				IsRegexp:    route.IsRegexp,
			})

			if result {
				return &MatchedRoute{
					Route:  &r.Routes[i],
					Params: params,
				}
			}
		}
	}

	return nil
}

func (r *Routing) Execute(route *MatchedRoute, request *Request, controller *Controller) {
	var logger = log.New(log.Writer(), "Routing.Execute", log.Flags())

	logger.Printf("Got request: %+v\n", request)

	defer func() {
		recovered := recover()

		if recovered != nil {
			r.Catch(recovered, controller, request, controller.response)
		}
	}()

	if route.Route.ContentType != "" {
		controller.response.Header().Add(HEADER_KEY_CONTENT_TYPE, route.Route.ContentType)
	}

	err := route.Route.Handler(request, controller)

	if err != nil {
		logger.Printf("Got error while execute request handler: %+v\n", err)
		r.Catch(err, controller, request, controller.response)
		return
	}

	logger.Printf("Got status: %d, response: %+v", controller.status, string(controller.content))
}

func (r *Routing) Catch(err any, controller *Controller, request *Request, response http.ResponseWriter) {
	var logger = log.New(log.Writer(), "Routing.Catch", log.Flags())

	response.Header()
	response.WriteHeader(http.StatusInternalServerError)
	response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	logger.Printf("Got error while handling request: %+v, response: %+v, error: %+v\n", request, string(controller.content), err)
}
