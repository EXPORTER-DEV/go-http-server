package server

import (
	"log"
	"net/http"
	"strconv"
)

type ServerInterface interface {
	Get(routes ...Route)
	Init() error
}

const (
	POST                    = "POST"
	GET                     = "GET"
	HEADER_KEY_CONTENT_TYPE = "Content-Type"
)

type Server struct {
	Port   int
	Host   string
	routes []Route
	server *http.Server
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
	Routing *Routing
}

func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	route := h.Routing.Match(request.Method, request.URL.Path)

	h.Routing.Execute(route, request, response)
}

func (s *Server) Init() error {
	var addr string = ""

	if s.Host != "" {
		addr = s.Host
	}

	if s.Port != 0 {
		addr = addr + ":" + strconv.Itoa(s.Port)
	}

	log.Printf("Starting server on %+v", addr)

	mux := http.NewServeMux()

	routing := Routing{
		Routes: s.routes,
	}

	mux.Handle("/", Handler{
		Routing: &routing,
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
