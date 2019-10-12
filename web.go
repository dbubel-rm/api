package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
	"time"
	log "github.com/sirupsen/logrus"
)

type MiddleWare func( Handler) Handler
type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params)
type App struct {
	Router           *httprouter.Router
	globalMiddleware []MiddleWare
}

func New() *App {
	//l.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	return &App{
		Router:           httprouter.New(),
	}
}

func (a *App) GlobalMiddleware(mid ...MiddleWare) {
	a.globalMiddleware = mid
}

func (a *App) SetLoggingLevel(level log.Level) {
	log.SetLevel(level)
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

	a.globalMiddleware = append([]MiddleWare{
		// Add a start timer middleware to the beginning of global middleware slice
		func(next Handler) Handler {
			return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
				ctx := context.WithValue(r.Context(), "ts", time.Now())
				next(w, r.WithContext(ctx), params)
			}
		},
	}, a.globalMiddleware...)

	// Wrap handler in global middleware
	for i := len(a.globalMiddleware) - 1; i >= 0; i-- {
		if a.globalMiddleware[i] != nil {
			finalHandler = a.globalMiddleware[i](finalHandler)
		}
	}

	h1 := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		// Essentially wrap our handler chain in a httprouter handle
		finalHandler(w, r, params)
	}

	a.Router.Handle(verb, path, h1)
	log.WithFields(log.Fields{"path": path}).Debug("added route")
}
