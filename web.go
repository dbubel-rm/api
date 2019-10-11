package api

import (
	"context"
	"github.com/julienschmidt/httprouter"

	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

//type rr map[string]interface{}
type MiddleWare func(handler http.Handler) http.Handler
type Handler func(w http.ResponseWriter, r *http.Request)

var l *log.Logger

type App struct {
	Router           *httprouter.Router
	globalMiddleware []MiddleWare
	logging          bool
}

// New creates an App value that handle a set of routes for the application.
func New(mux *httprouter.Router, mw ...MiddleWare) *App {
	l = log.New()
	//l.SetReportCaller(true)
	l.SetFormatter(&log.JSONFormatter{})
	l.SetOutput(os.Stdout)
	l.SetLevel(log.DebugLevel)

	return &App{
		Router:           mux,
		globalMiddleware: mw,
	}
}

func (a *App) SetLoggingLevel(level log.Level) {
	l.SetLevel(level)
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

	builtinMiddlwares := []MiddleWare{
		// Add a start timer middleware to the beginning of global middlware slice
		func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := context.WithValue(r.Context(), "ts", time.Now())
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		},
	}

	a.globalMiddleware = append(builtinMiddlwares, a.globalMiddleware...)

	// Wrap handler in global middleware
	for i := len(a.globalMiddleware) - 1; i >= 0; i-- {
		if a.globalMiddleware[i] != nil {
			h = a.globalMiddleware[i](h)
		}
	}

	l.WithFields(log.Fields{"path": path}).Info("added route")
	a.Router.Handler(verb, path, h)
}
