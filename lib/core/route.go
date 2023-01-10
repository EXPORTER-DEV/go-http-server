package server

import (
	"log"
	"net/http"
	"regexp"
)

type RouteHandler func(request *Request, controller *Controller) error

type Route struct {
	Method      string
	Handler     RouteHandler
	ContentType string
	Path        string
	IsRegexp    bool
	HasParams   bool
}

type RouteOptions struct {
	HasParams bool
	IsRegexp  bool
}

func NewRoute(path string, handler RouteHandler, options ...RouteOptions) *Route {
	routeOptions := RouteOptions{}

	if len(options) > 0 {
		routeOptions = options[0]
	}

	return &Route{
		Handler:   handler,
		Path:      path,
		IsRegexp:  routeOptions.IsRegexp,
		HasParams: routeOptions.HasParams,
	}
}

type Params map[string]string

type MatchedRoute struct {
	Route  *Route
	Params Params
}

type Routing struct {
	Routes []Route
}

func (r *Routing) Match(method string, path string) *MatchedRoute {
	for i, route := range r.Routes {
		var isEqual bool = false
		var params = make(map[string]string)

		if route.Method == method {
			if route.IsRegexp {
				r := regexp.MustCompile(route.Path)
				if r.Match([]byte(path)) {
					isEqual = true
				}
			}
			if route.HasParams {
				routeParamsRegexp := regexp.MustCompile(`(\:[a-zA-Z]+)`)

				routeParamsMatch := routeParamsRegexp.FindAllSubmatchIndex([]byte(route.Path), -1)

				if len(routeParamsMatch) > 0 {
					routePathRegexp := route.Path

					routeParamNames := []string{}

					for _, indexes := range routeParamsMatch {
						name := routePathRegexp[indexes[2]:indexes[3]]

						routeParamNames = append(routeParamNames, name)

						routePathRegexp = routePathRegexp[:indexes[0]] + "(.*)" + routePathRegexp[indexes[1]:]
					}

					pathRegexp := regexp.MustCompile(routePathRegexp)

					pathMatch := pathRegexp.FindAllSubmatchIndex([]byte(path), -1)

					if len(pathMatch) > 0 {
						params := make(map[string]string)

						for index, matchIndexes := range pathMatch {
							value := path[matchIndexes[2]:matchIndexes[3]]

							params[routeParamNames[index]] = value
						}

						isEqual = true
					}
				}
			}
			if !isEqual && route.Path == path {
				isEqual = true
			}

			if isEqual {
				return &MatchedRoute{
					Route:  &r.Routes[i],
					Params: params,
				}
			}
		}
	}

	return nil
}

func (r *Routing) Execute(route *MatchedRoute, request *http.Request, response http.ResponseWriter) {
	var logger = log.New(log.Writer(), "Routing.Execute", log.Flags())

	if route == nil {
		logger.Printf("Got no handler for request (%v) %+v\n", request.Method, request.URL.Path)
		return
	}

	var controller = NewController(request, response)
	var requestWrapper, err = NewRequest(request, route.Params)

	logger.Printf("Got request: %+v\n", requestWrapper)

	if err != nil {
		logger.Printf("Got error while prepare request: %+v\n", err)
		r.Catch(err, controller, requestWrapper, response)
		return
	}

	defer func() {
		recovered := recover()

		if recovered != nil {
			r.Catch(recovered, controller, requestWrapper, response)
		}
	}()

	err = route.Route.Handler(requestWrapper, controller)

	if err != nil {
		logger.Printf("Got error while execute request handler: %+v\n", err)
		r.Catch(err, controller, requestWrapper, response)
		return
	}

	contentType := response.Header().Get(HEADER_KEY_CONTENT_TYPE)

	if contentType == "" && route.Route.ContentType != "" {
		response.Header().Add(HEADER_KEY_CONTENT_TYPE, route.Route.ContentType)
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
