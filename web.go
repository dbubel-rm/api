package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type rr map[string]interface{}
type MiddleWare func(handler http.Handler) http.Handler

type App struct {
	Router           *http.ServeMux
	log              *log.Logger
	globalMiddleware []MiddleWare
}

// New creates an App value that handle a set of routes for the application.
func New(mux *http.ServeMux, mw ...MiddleWare) *App {
	return &App{
		Router:           mux,
		log:              log.New(os.Stdout, "", 0),
		globalMiddleware: mw,
	}
}

func GET(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
		} else {
			Respond(next, w, r, rr{"error": "Method Now Allowed"}, http.StatusMethodNotAllowed)
		}
	})
}

func POST(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		} else {
			Respond(next, w, r, rr{"error": "Method Now Allowed"}, http.StatusMethodNotAllowed)
		}
	})
}

func Start(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), "ts", time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func End(l *log.Logger) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		responseTime := time.Now().Sub(r.Context().Value("ts").(time.Time))
		fmt.Println(responseTime.String())
		json, _ := json.Marshal(rr{
			"method":       r.Method,
			"url":          r.RequestURI,
			"contentLenth": r.ContentLength,
			"ip":           r.RemoteAddr,
			"ts":           responseTime.String(),
		})
		l.Println(string(json))
	}
}

func (a *App) Handle(verb MiddleWare, path string, finalHandler MiddleWare, middlwares ...MiddleWare) {

	var h http.Handler
	h = http.HandlerFunc(End(a.log))

	// Route middle wares
	middlwares = append(middlwares, finalHandler)

	for i := len(middlwares) - 1; i >= 0; i-- {
		if middlwares[i] != nil {
			h = middlwares[i](h)
		}
	}

	a.globalMiddleware = append(a.globalMiddleware, verb)
	//if verb == http.MethodGet {
	//	a.globalMiddleware = append(a.globalMiddleware, GET)
	//} else if verb == http.MethodPost {
	//	a.globalMiddleware = append(a.globalMiddleware, POST)
	//} else {
	//	log.Fatalln("Invalid Method")
	//}

	a.globalMiddleware = append([]MiddleWare{Start}, a.globalMiddleware...)
	for i := len(a.globalMiddleware) - 1; i >= 0; i-- {
		if a.globalMiddleware[i] != nil {
			h = a.globalMiddleware[i](h)
		}
	}

	a.log.Println("Adding route for", path)
	a.Router.Handle(path, h)
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

