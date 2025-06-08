package server

import (
	"context"
)

type Application struct {
}

func New() *Application {

	return &Application{}
}

func (a *Application) Start(ctx context.Context) {

}
