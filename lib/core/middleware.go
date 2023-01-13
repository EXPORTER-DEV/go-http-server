package server

import "exporter-dev/http-server/lib/utils"

type Middleware struct {
	Method      []string
	Handler     MiddlewareHandler
	Path        []string
	ParseParams bool
	IsRegexp    bool
}

type MiddlewareHandler func(request *Request, controller *Controller) (skip bool, err error)

func NewMiddleware(handler MiddlewareHandler) *Middleware {
	return &Middleware{
		Handler: handler,
	}
}

func (m *Middleware) SetIsRegexp(value bool) *Middleware {
	m.IsRegexp = value

	return m
}

func (m *Middleware) SetParseParams(value bool) *Middleware {
	m.ParseParams = value

	return m
}

func (m *Middleware) AddPath(path string) *Middleware {
	m.Path = append(m.Path, path)

	return m
}

type Middlewaring struct {
	Middlewares []Middleware
}

func (m *Middlewaring) Execute(request *Request, controller *Controller) (skip bool, err error) {
	method := request.Method
	path := request.Path

	var middlewares []Middleware

	for index, middleware := range m.Middlewares {
		if len(middleware.Method) > 0 &&
			!utils.Some(middleware.Method, func(item string, index int) bool {
				return item == method
			}) {
			continue
		}

		if len(middleware.Method) == 0 && len(middleware.Path) == 0 {
			middlewares = append(middlewares, middleware)

			continue
		}

		var validate bool

		if utils.Some(middleware.Path, func(item string, index int) bool {
			return item == request.Route.Path
		}) {
			validate = true
		}

		if !validate {
		inner:
			for _, handlerPath := range middleware.Path {
				matching := &Matching{
					RequestedPath: path,
					HandlerPath:   handlerPath,
				}

				result, _ := matching.Execute(
					MatchingExecuteOptions{
						IsRegexp:    middleware.IsRegexp,
						ParseParams: middleware.ParseParams,
					},
				)

				if result {
					validate = true
					break inner
				}
			}
		}

		if validate {
			middlewares = append(middlewares, middlewares[index])
		}
	}

	return m.handle(middlewares, request, controller)
}

type ConcurrencyMiddlewareResult struct {
	skip bool
	err  error
}

func handleConcurrencyMiddleware(middleware *Middleware, request *Request, controller *Controller, output chan<- ConcurrencyMiddlewareResult) {
	skip, err := middleware.Handler(request, controller)

	result := ConcurrencyMiddlewareResult{
		skip,
		err,
	}

	output <- result
}

func handleMiddleware(middleware *Middleware, request *Request, controller *Controller) (skip bool, err error) {
	output := make(chan ConcurrencyMiddlewareResult)
	go handleConcurrencyMiddleware(middleware, request, controller, output)

	result := <-output

	return result.skip, result.err
}

func (m *Middlewaring) handle(middlewares []Middleware, request *Request, controller *Controller) (skip bool, err error) {
	if len(middlewares) > 0 {
		current := middlewares[0]

		skip, err := handleMiddleware(&current, request, controller)

		if skip {
			return true, nil
		}

		if err == nil {
			return m.handle(middlewares[1:], request, controller)
		}

		return false, err
	}

	return false, nil
}
