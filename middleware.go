package tgbotapp

type Middleware func(context *BotContext, next HandlerFunc)

type MiddlewareChain struct {
	middlewares []Middleware
}

func NewMiddlewareChain() *MiddlewareChain {
	return &MiddlewareChain{
		middlewares: make([]Middleware, 0),
	}
}

func (c *MiddlewareChain) Append(middleware ...Middleware) {

	c.middlewares = append(c.middlewares, middleware...)
}

func (c *MiddlewareChain) Wrap(final HandlerFunc) HandlerFunc {

	return func(ctx *BotContext) {

		var exec func(int, *BotContext)

		exec = func(index int, ctx *BotContext) {
			if index < len(c.middlewares) {
				c.middlewares[index](ctx, func(ctx1 *BotContext) {
					exec(index+1, ctx1)
				})
			} else {
				final(ctx)
			}
		}

		exec(0, ctx)
	}

}
