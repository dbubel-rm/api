package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
)

var apiLogger *logrus.Logger

type MiddleWare func(Handler) Handler
type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params)

type App struct {
	Router            *httprouter.Router
	globalMiddlewares []MiddleWare
}

func New(log *logrus.Logger) *App {
	apiLogger = log
	return &App{
		Router:            httprouter.New(),
		globalMiddlewares: make([]MiddleWare, 0, 0),
	}
}

func NewDefault() *App {
	apiLogger = logrus.New()
	apiLogger.SetLevel(logrus.DebugLevel)
	apiLogger.SetFormatter(&logrus.JSONFormatter{})
	return &App{
		Router:            httprouter.New(),
		globalMiddlewares: make([]MiddleWare, 0, 0),
	}
}

func (a *App) GlobalMiddleware(mid ...MiddleWare) {
	a.globalMiddlewares = mid
}

func (a *App) Endpoints(e ...Endpoints) {
	for x := 0; x < len(e); x++ {
		for i := 0; i < len(e[x]); i++ {
			a.Handle(e[x][i].Verb, e[x][i].Path, e[x][i].EndpointHandler, e[x][i].MiddlewareHandlers...)
		}
	}
}

func (a *App) Handle(verb string, path string, finalHandler Handler, middleware ...MiddleWare) {
	// Wrap all the route specific middleware
	for i := len(middleware) - 1; i >= 0; i-- {
		if middleware[i] != nil {
			finalHandler = middleware[i](finalHandler)
		}
	}

	a.globalMiddlewares = append([]MiddleWare{
		// Add a start timer middleware to the beginning of global middleware slice
		func(next Handler) Handler {
			return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
				ctx := context.WithValue(r.Context(), "ts", time.Now())
				next(w, r.WithContext(ctx), params)
			}
		},
	}, a.globalMiddlewares...)

	// Wrap handler in global middleware
	for i := len(a.globalMiddlewares) - 1; i >= 0; i-- {
		if a.globalMiddlewares[i] != nil {
			finalHandler = a.globalMiddlewares[i](finalHandler)
		}
	}

	// Our wrapped function chain in a compatible httprouter Handle func
	a.Router.Handle(verb, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		finalHandler(w, r, params)
	})
	apiLogger.WithFields(logrus.Fields{"path": path}).Debug("added route")
}

func (a *App) Run(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	apiLogger.WithFields(logrus.Fields{"addr": server.Addr}).Info("server starting")

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		apiLogger.WithError(err).Error("error starting server")
	case <-osSignals:
		apiLogger.Info("shutdown received shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			apiLogger.WithError(err).Error("graceful shutdown did not complete in allowed time")
			if err := server.Close(); err != nil {
				apiLogger.WithError(err).Error("could not stop http server")
			}
		}
		apiLogger.Info("shutdown OK")
	}
}
