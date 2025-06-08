package server

import (
	tgbotapp "go-telegram-bot-app/v1"
	"log"
)

type Context struct {
	logger *log.Logger
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

func (c *Context) Logger() *log.Logger {

	if c.logger == nil {
		panic(tgbotapp.ErrLoggerNotFound)
	}

	return c.logger
}

func (c *Context) SetLogger(logger *log.Logger) {
	c.logger = logger
}
