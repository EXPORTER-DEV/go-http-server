package server

type Context struct {
	Store map[string]string
}

func (c *Context) Get(key string) string {
	return c.Store[key]
}

func (c *Context) Set(key string, value string) {
	c.Store[key] = value
}
