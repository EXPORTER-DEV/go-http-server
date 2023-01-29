# go-http-server

Custom HTTP Server implementation with Routing (with RegExp, simple path, params support), custom Request (with Context and Binding it), Controller implementations and Middleware based on `net/http` lib.

*To run tests for `server` package:*

```bash
$ go test ./lib/core -v
```

### Usage

*Example of register Server instance:*

```go

import (
	server "exporter-dev/http-server/lib/core"
)

func main() {
	// Create new Server instance, and SetPort (by default will use 80 port) and SetHost:
	var instance = server.NewServer().SetPort(3000).SetHost("domain.com");

    // Register get handler for route for path "/index":
    instance.Get(*server.NewRoute("/index", func(request *server.Request, controller *server.Controller) error {
		controller.Status(200)
		controller.Send("Hello, world!")
		return nil
	}))

    // Init server listening, after init we can't SetPort or SetHost for the server instance:
    instance.init()
}
```

**Examples of usage middlewares, context binding, and route group:**

Check the source code example at *[./main.go](./main.go)*

### Todo

- [ ] Increase unit tests cover

- [ ] Optimize route matching when ParseParams is true and every matching request attemp recompiles the RegExp for the route

- [ ] Make default response status code for each Route