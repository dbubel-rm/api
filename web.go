package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type MiddleWare func(Handler) Handler
type Handler func(w http.ResponseWriter, r *http.Request, params httprouter.Params)
var apiLogger *logrus.Logger

type App struct {
	Router            *httprouter.Router
	globalMiddlewares []MiddleWare
}

func NewBasic(logger *logrus.Logger) *App {
	apiLogger = logger
	return &App{
		Router: httprouter.New(),
	}
}

func (a *App) GlobalMiddleware(mid ...MiddleWare) {
	a.globalMiddlewares = mid
}

func (a *App) Endpoints(e Endpoints) {
	for i := 0; i < len(e); i++ {
		a.Handle(e[i].Method, e[i].Path, e[i].EndpointHandler, e[i].MiddlewareHandlers...)
	}
}

func (a *App) Handle(verb string, path string, finalHandler Handler, middlwares ...MiddleWare) {
	// Wrap all the route specific middleware
	for i := len(middlwares) - 1; i >= 0; i-- {
		if middlwares[i] != nil {
			finalHandler = middlwares[i](finalHandler)
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

	a.Router.Handle(verb, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		finalHandler(w, r, params) // our wrapped function chain
	})
	apiLogger.WithFields(logrus.Fields{"path": path}).Debug("added route")
}

func(a *App) StartAPI(addr string) {
	server := http.Server{
		Addr:           addr,
		Handler:        a.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		apiLogger.WithError(err).Error("Error starting server")
	case <-osSignals:
		apiLogger.Info("shutdown signal recieved shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second * 11)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			apiLogger.WithError(err).Error("Graceful shutdown did not complete in allowed time")
			if err := server.Close(); err != nil {
				apiLogger.WithError(err).Error("Could not stop http server")
			}
		}
		apiLogger.Info("Shutdown OK")
	}
}
