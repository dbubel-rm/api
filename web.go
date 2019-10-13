package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type MiddleWare func(Handler) Handler
type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params)
type App struct {
	Router         *httprouter.Router
	appMiddlewares []MiddleWare
	log            *logrus.Logger
}

func NewBasic(logger *logrus.Logger) *App {
	return &App{
		Router: httprouter.New(),
		log:    logger,
	}
}

func (a *App) GlobalMiddleware(mid ...MiddleWare) {
	a.appMiddlewares = mid
}

func (a *App) Endpoints(endpoints Endpoints) {
	for i := 0; i < len(endpoints); i++ {
		a.Handle(endpoints[i].Method, endpoints[i].Path, endpoints[i].EndpointHandler, endpoints[i].MiddlewareHandlers...)
	}
}

func (a *App) Handle(verb string, path string, finalHandler Handler, middlwares ...MiddleWare) {
	// Wrap all the route specific middleware
	for i := len(middlwares) - 1; i >= 0; i-- {
		if middlwares[i] != nil {
			finalHandler = middlwares[i](finalHandler)
		}
	}

	a.appMiddlewares = append([]MiddleWare{
		// Add a start timer middleware to the beginning of global middleware slice
		func(next Handler) Handler {
			return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
				ctx := context.WithValue(r.Context(), "ts", time.Now())
				next(w, r.WithContext(ctx), params)
			}
		},
	}, a.appMiddlewares...)

	// Wrap handler in global middleware
	for i := len(a.appMiddlewares) - 1; i >= 0; i-- {
		if a.appMiddlewares[i] != nil {
			finalHandler = a.appMiddlewares[i](finalHandler)
		}
	}

	a.Router.Handle(verb, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		finalHandler(w, r, params) // our wrapped function chain
	})
	a.log.WithFields(logrus.Fields{"path": path}).Debug("added route")
}
