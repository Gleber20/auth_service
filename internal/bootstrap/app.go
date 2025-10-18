package bootstrap

import (
	"auth_service/internal/config"
	"context"
	"errors"
	"net/http"
	"time"
)

const gracefulDeadline = 5 * time.Second

type App struct {
	cfg      config.Config
	rest     *http.Server
	teardown []func()
}

func New(cfg config.Config) *App {
	app := initLayers(cfg)

	return app
}

func (app *App) Run(ctx context.Context) {
	go func() {
		if err := app.rest.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done()

	for i := range app.teardown {
		app.teardown[i]()
	}
}

func (app *App) HTTPHandler() http.Handler {
	return app.rest.Handler
}
