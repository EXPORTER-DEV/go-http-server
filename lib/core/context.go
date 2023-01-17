package server

import "reflect"

type Context struct {
	Store map[string]any
}

func (c *Context) Get(key string) (any, bool) {
	value, exists := c.Store[key]

	return value, exists
}

func (c *Context) Set(key string, value any) {
	c.Store[key] = value
}

func BindContext[T any](context *Context, key string) T {
	value, exists := context.Get(key)

	if !exists {
		return *new(T)
	}

	v := reflect.ValueOf(value)

	if v.IsZero() {
		return *new(T)
	}

	return v.Interface().(T)
}
