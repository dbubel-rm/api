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

//func (a *App) allowMethod(method string) func(next http.Handler) http.Handler {
//	return func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			if r.Method == method {
//				next.ServeHTTP(w, r)
//			} else {
//				Respond(w, r, map[string]interface{}{"error": "Method Now Allowed"}, http.StatusMethodNotAllowed)
//			}
//		})
//	}
//}

func (a *App) Handle(verb string, path string, finalHandler Handler, middlwares ...MiddleWare) {

	// Our handler that we will pass to the router
	var h http.Handler = http.HandlerFunc(finalHandler)

	// Wrap all the route specific middleware
	for i := len(middlwares) - 1; i >= 0; i-- {
		if middlwares[i] != nil {
			h = middlwares[i](h)
		}
	}

	// Assign a middleware that will only allow the specified verb
	//if verb == http.MethodGet {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodGet))
	//} else if verb == http.MethodOptions {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodOptions))
	//} else if verb == http.MethodDelete {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodDelete))
	//} else if verb == http.MethodPost {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodPost))
	//} else if verb == http.MethodPut {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodPut))
	//} else if verb == http.MethodConnect {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodConnect))
	//} else if verb == http.MethodHead {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodHead))
	//} else if verb == http.MethodTrace {
	//	a.globalMiddleware = append(a.globalMiddleware, a.allowMethod(http.MethodTrace))
	//}

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

	//a.debug(rr{"path": path, "msg": "added route"})
	l.WithFields(log.Fields{"path":path}).Info("added route")
	//a.Router.Handle(path, h)
	a.Router.Handler(verb,path,h)
	//a.Router.HandleFunc()

	//finalHandler := http.HandlerFunc(final)
	//x := middlewareOne(middlewareTwo(finalHandler))
	//http.Handle("/", h)
	// Wrap up the application-wide first, this will call the first function
	// of each middleware which will return a function of type Handler.``
	//finalHandler = wrapMiddleware(finalHandler, middlwares)
	//finalHandler = wrapMiddleware(finalHandler, a.middlwares)

	// The function to execute for each request.
	//h := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Setup CORS for development - Dont do this unless its a public endpoint obviously
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	// w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// Call the wrapped finalHandler functions.
	//a.Router.ServeHTTP(w, r)
	//finalHandler(a.log, w, r)
	//}
	// Add this finalHandler for the specified verb and route.
	//a.Router.Handle(verb, path, h)
	//a.Router.HandleFunc(path, h)
	//a.Handle(verb, path, h)

	//Handle preflight browser requests when we are POSTing data
	// f := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// 	// Setup CORS for dev only or public endpoint
	// 	w.Header().Set("Access-Control-Allow-Origin", "*")
	// 	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	// 	return
	// }
	// a.Router.Handle("OPTIONS", path, f)
}
