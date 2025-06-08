package middleware

import (
	"go-telegram-bot-app/v1/internal/server"
)

type HandlerFunc func(context *server.Context)

type Middleware func(context *server.Context, next HandlerFunc)

type MiddlewareChain struct {
	middlewares []Middleware
}

func New() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

func (c *MiddlewareChain) Append(middleware ...Middleware) {

	c.middlewares = append(c.middlewares, middleware...)
}

func (c *MiddlewareChain) Wrap(final HandlerFunc) HandlerFunc {

	return func(ctx *server.Context) {

		var exec func(int, *server.Context)

		exec = func(index int, ctx *server.Context) {
			if index < len(c.middlewares) {
				c.middlewares[index](ctx, func(ctx1 *server.Context) {
					exec(index+1, ctx1)
				})
			} else {
				final(ctx)
			}
		}

		exec(0, ctx)
	}

}
