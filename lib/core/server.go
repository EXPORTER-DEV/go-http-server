package server

import (
	"io"
	"log"
	"net/http"
	"strconv"
)

type ServerHandlerFunc func(request *http.Request, response http.ResponseWriter) error

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

func (h Handler) ExecuteRoutes(routes []*Route, request *http.Request, response http.ResponseWriter) error {
	var controller = NewController(request, response)

	for _, route := range routes {
		err := route.Handler(request, controller)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h Handler) MatchRoutes(method string, path string) []*Route {
	var result []*Route

	for _, route := range *h.Routes {
		if route.Method == method && route.Path == path {
			result = append(result, &route)
		}
	}

	return result
}

func (h Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	defer func() {
		r := recover()

		if r != nil {
			response.Header()
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte(http.StatusText(http.StatusInternalServerError)))
			log.Printf("Got error while handling request: %+v, response: %+v\n", r, response)
			return
		}

		contentType := response.Header().Get(HEADER_KEY_CONTENT_TYPE)

		if contentType == "" && *h.ContentType != "" {
			response.Header().Add(HEADER_KEY_CONTENT_TYPE, *h.ContentType)
		}

		log.Printf("Got response for request: %+v", response)
	}()

	defer request.Body.Close()

	routes := h.MatchRoutes(request.Method, request.URL.Path)

	bodyBytes, err := io.ReadAll(request.Body)

	body := string(bodyBytes)

	if err != nil {
		log.Printf("Got error while parsing request body: %+v\n", err)
		panic("FAILED_PARSE_REQUEST_BODY")
	}

	requestForLog := RequestForLog{
		Method:     request.Method,
		Url:        request.URL.String(),
		Body:       body,
		RemoteAddr: request.RemoteAddr,
	}

	log.Printf("Got incoming request: %+v\n", requestForLog)

	if len(routes) == 0 {
		log.Printf("Not found handler for request")
		panic("NOT_FOUND_HANDLER")
	}

	err = h.ExecuteRoutes(routes, request, response)

	if err != nil {
		log.Printf("Failed execution handler for request: %+v", err)
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

func (s *Server) Get(routes ...Route) {
	for index := range routes {
		routes[index].Method = "GET"
	}

	s.routes = append(s.routes, routes...)
}
