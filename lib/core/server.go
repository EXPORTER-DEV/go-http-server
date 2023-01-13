package server

import (
	"log"
	"net/http"
	"strconv"
)

const (
	POST                    = "POST"
	GET                     = "GET"
	HEADER_KEY_CONTENT_TYPE = "Content-Type"
)

type Server struct {
	Port        int
	Host        string
	routes      []Route
	middlewares []Middleware
	server      *http.Server
}

type RequestForLog struct {
	Url        string
	Method     string
	Body       string
	RemoteAddr string
}

type ResponseForLog struct {
	Body   string
	Status int
	Header http.Header
}

type Handler struct {
	Routing      *Routing
	Middlewaring *Middlewaring
}

func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	var logger = log.New(log.Writer(), "Routing.Execute", log.Flags())

	route := h.Routing.Match(request.Method, request.URL.Path)

	if route == nil {
		logger.Printf("Got no handler for request (%v) %+v\n", request.Method, request.URL.Path)
		return
	}

	var controller = NewController(request, response)
	var requestWrapper, err = NewRequest(request, route.Params, route.Route)

	if err != nil {
		logger.Printf("Failed prepare request: %+v, error: %+v\n", request, err)
		h.Routing.Catch(err, controller, requestWrapper, response)
		return
	}

	logger.Printf("Got request: %+v\n", requestWrapper)

	skip, err := h.Middlewaring.Execute(requestWrapper, controller)

	if err != nil {
		logger.Printf("Got error while handling middlewares: %+v", err)
		h.Routing.Catch(err, controller, requestWrapper, response)
		return
	}

	if !skip {
		h.Routing.Execute(route, requestWrapper, controller)
	}
}

func (s *Server) Init() error {
	var addr string = ""

	if s.Host != "" {
		addr = s.Host
	}

	if s.Port != 0 {
		addr = addr + ":" + strconv.Itoa(s.Port)
	}

	log.Printf("Starting server on %+v, with routes (%d) & middlewares (%d)", addr, len(s.routes), len(s.middlewares))

	mux := http.NewServeMux()

	mux.Handle("/", Handler{
		Routing: &Routing{
			Routes: s.routes,
		},
		Middlewaring: &Middlewaring{
			Middlewares: s.middlewares,
		},
	})

	s.server = &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	error := s.server.ListenAndServe()

	if error != nil {
		log.Fatalf("Got error while starting server: %s", error)
		return error
	}
	return nil
}

func (s *Server) Get(routes ...Route) {
	for index := range routes {
		routes[index].Method = "GET"
	}

	s.routes = append(s.routes, routes...)
}

func (s *Server) Use(middlewares ...Middleware) {
	s.middlewares = append(s.middlewares, middlewares...)
}
