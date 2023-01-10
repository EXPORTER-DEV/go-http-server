package server

import (
	"io"
	"net/http"
	"strings"
)

type ControllerHeader struct {
	controller *Controller
}

type ControllerInterface interface {
	Send(content string) error
	Append(bytes []byte)
	Status(status int)
	Header() *ControllerHeader
}

type Response struct {
	Content *[]byte
	Status  *int
	Headers map[string][]string
}

type Request struct {
	Headers    map[string][]string
	Body       string
	Content    *[]byte
	Query      map[string][]string
	Path       string
	Url        string
	Params     Params
	RemoteAddt string
	Context    interface{}
}

func getRequestBody(request *http.Request) (content *[]byte, body string, err error) {
	bytes, err := io.ReadAll(request.Body)

	defer func() {
		request.Body.Close()
	}()

	if err != nil {
		return &[]byte{}, "", err
	}

	content = &bytes

	body = string(bytes)

	return
}

func NewRequest(request *http.Request, params Params) (*Request, error) {
	content, body, err := getRequestBody(request)

	if err != nil {
		return nil, err
	}

	var context interface{}

	return &Request{
		Headers:    request.Header,
		Body:       body,
		Content:    content,
		Query:      request.URL.Query(),
		Path:       request.URL.Path,
		Url:        request.URL.String(),
		Params:     params,
		RemoteAddt: request.RemoteAddr,
		Context:    context,
	}, nil
}

type Controller struct {
	request  *http.Request
	response http.ResponseWriter
	headers  map[string][]string
	status   int
	content  []byte
	Header   *ControllerHeader
	Response *Response
}

func NewController(request *http.Request, response http.ResponseWriter) *Controller {
	var controller *Controller = &Controller{
		request:  request,
		response: response,
		status:   200,
		headers:  make(map[string][]string),
	}

	controller.Header = &ControllerHeader{
		controller,
	}

	controller.Response = &Response{
		Content: &controller.content,
		Status:  &controller.status,
		Headers: controller.headers,
	}

	return controller
}

// Send string content to client
func (controller *Controller) Send(content string) error {
	if len(content) > 0 {
		controller.content = append(controller.content, []byte(content)...)
	}

	for key, values := range controller.headers {
		for _, value := range values {
			controller.response.Header().Add(key, value)
		}
	}

	controller.response.WriteHeader(controller.status)

	_, err := controller.response.Write(controller.content)

	return err
}

// Append bytes of response buffer, not sending to client
func (controller *Controller) Append(bytes []byte) {
	controller.content = append(controller.content, bytes...)
}

// Set status of response, not sending to client
func (controller *Controller) Status(status int) {
	controller.status = status
}

// Get header by key, if ignoreCase is true - will ignore case when matching headers
func (header *ControllerHeader) Get(key string, ignoreCase bool) []string {
	var found []string = []string{}

	if ignoreCase {
		key = strings.ToLower(key)
	}

	for headerKey, values := range header.controller.headers {
		if ignoreCase {
			headerKey = strings.ToLower(headerKey)
		}

		if headerKey == key {
			found = append(found, values...)
		}
	}

	return found
}

// Clear all set headers
func (header *ControllerHeader) Clear() {
	header.controller.headers = map[string][]string{}
}

// Remove headers by key
func (header *ControllerHeader) Remove(key string, ignoreCase bool) {
	if ignoreCase {
		key = strings.ToLower(key)
	}

	delete(header.controller.headers, key)
}

// Add header by key and value
func (header *ControllerHeader) Add(key string, value string) {
	if header.controller.headers[key] != nil {
		header.controller.headers[key] = append(header.controller.headers[key], value)
	} else {
		header.controller.headers[key] = []string{value}
	}
}
