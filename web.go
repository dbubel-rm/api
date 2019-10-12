package api

import (
	"context"

	"github.com/julienschmidt/httprouter"

	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type MiddleWare func(handler http.Handler) http.Handler
type Handler func(w http.ResponseWriter, r *http.Request)

type App struct {
	Router           *httprouter.Router
	globalMiddleware []MiddleWare
}

// New creates an App value that handle a set of routes for the application.
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
	// Our handler that we will pass to the router
	var h http.Handler = http.HandlerFunc(finalHandler)

	// Wrap all the route specific middleware
	for i := len(middlwares) - 1; i >= 0; i-- {
		if middlwares[i] != nil {
			h = middlwares[i](h)
		}
	}

	a.globalMiddleware = append([]MiddleWare{
		// Add a start timer middleware to the beginning of global middleware slice
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "ts", time.Now())
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}, a.globalMiddleware...)

	// Wrap handler in global middleware
	for i := len(a.globalMiddleware) - 1; i >= 0; i-- {
		if a.globalMiddleware[i] != nil {
			h = a.globalMiddleware[i](h)
		}
	}

	log.WithFields(log.Fields{"path": path}).Debug("added route")
	a.Router.Handler(verb, path, h)
}
