package server

import (
	"context"
	"log"
)

type Context struct {
	Logger *log.Logger
	data   map[string]any
}

func NewContext() *Context {
	return &Context{

		data: make(map[string]any),
	}
}

func (c *Context) SetData(key string, value any) {
	c.data[key] = value
}

func (c *Context) GetData(key string) (value any, ok bool) {
	value, ok = c.data[key]
	return
}

type Application struct {
}

func New() *Application {

	return &Application{}
}

func (a *Application) Start(ctx context.Context) {

}
