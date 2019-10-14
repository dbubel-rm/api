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

type App struct {
	Router            *httprouter.Router
	globalMiddlewares []MiddleWare
}

func New() *App {
	ApiLogger = logrus.New()
	ApiLogger.SetFormatter(&logrus.JSONFormatter{})
	ApiLogger.SetLevel(logrus.DebugLevel)
	return &App{
		Router: httprouter.New(),
	}
}

func (a *App) GlobalMiddleware(mid ...MiddleWare) {
	a.globalMiddlewares = mid
}

func (a *App) Endpoints(ep ...Endpoints) {
	for x := 0; x < len(ep); x++ {
		for i := 0; i < len(ep[x]); i++ {
			a.Handle(ep[x][i].Method, ep[x][i].Path, ep[x][i].EndpointHandler, ep[x][i].MiddlewareHandlers...)
		}
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

	// Our wrapped function chain in a compatible httprouter Handle func
	a.Router.Handle(verb, path, func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		finalHandler(w, r, params)
	})
	ApiLogger.WithFields(logrus.Fields{"path": path}).Debug("added route")
}

func StartAPI(server *http.Server) {
	serverErrors := make(chan error, 1)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM)
	go func() {
		serverErrors <- server.ListenAndServe()
	}()

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		ApiLogger.WithError(err).Error("Error starting server")
	case <-osSignals:
		ApiLogger.Info("shutdown received shedding connections...")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*11)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			ApiLogger.WithError(err).Error("graceful shutdown did not complete in allowed time")
			if err := server.Close(); err != nil {
				ApiLogger.WithError(err).Error("could not stop http server")
			}
		}
		ApiLogger.Info("shutdown OK")
	}
}
