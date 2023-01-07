package server

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"golang.org/x/exp/slices"
)

type ServerHandlerFunc func(request *http.Request, response http.ResponseWriter) error

type ServerInterface interface {
	Get(path string, handler ServerHandlerFunc) error
	Init() error
}

const (
	POST                    = "POST"
	GET                     = "GET"
	HEADER_KEY_CONTENT_TYPE = "Content-Type"
)

type Route struct {
	method  string
	handler ServerHandlerFunc
	path    string
}

type Server struct {
	Port        int
	Host        string
	ContentType string
	routes      []Route
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
	Routes      *[]Route
	ContentType *string
}

func (h Handler) MatchRoute(method string, path string) (bool, *Route) {
	idx := slices.IndexFunc(*h.Routes, func(route Route) bool {
		return route.method == method && route.path == path
	})

	if idx == -1 {
		return false, &Route{}
	}

	route := &(*h.Routes)[idx]

	return true, route
}

func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer func() {
		r := recover()

		if r != nil {
			response.Header()
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			log.Printf("Got error while handling request, sent response: %+v\n", response)
			return
		}

		contentType := response.Header().Get(HEADER_KEY_CONTENT_TYPE)

		if contentType == "" && *h.ContentType != "" {
			response.Header().Add(HEADER_KEY_CONTENT_TYPE, *h.ContentType)
		}

		log.Printf("Got response for request: %+v", response)
	}()

	defer request.Body.Close()

	found, route := h.MatchRoute(request.Method, request.URL.Path)

	bodyBytes, bodyErr := io.ReadAll(request.Body)

	body := string(bodyBytes)

	if bodyErr != nil {
		log.Printf("Got error while parsing request body: %+v\n", bodyErr)
		panic("FAILED_PARSE_REQUEST_BODY")
	}

	requestForLog := RequestForLog{
		Method:     request.Method,
		Url:        request.URL.String(),
		Body:       body,
		RemoteAddr: request.RemoteAddr,
	}

	log.Printf("Got incoming request: %+v\n", requestForLog)

	if !found {
		log.Printf("Not found handler for request")
		panic("NOT_FOUND_HANDLER")
	}

	handlerErr := route.handler(request, response)

	if handlerErr != nil {
		log.Printf("Failed execution handler for request: %+v", handlerErr)
		panic("FAILED_EXECUTION_HANDLER")
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

	log.Printf("Starting server on %+v", addr)

	mux := http.NewServeMux()

	mux.Handle("/", Handler{
		Routes:      &s.routes,
		ContentType: &s.ContentType,
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

func (s *Server) Get(path string, handler ServerHandlerFunc) error {
	s.routes = append(s.routes, Route{
		path:    path,
		handler: handler,
		method:  GET,
	})

	return nil
}
