package web

import (
	"fmt"
	"log"
	"net/http"

)

//type Handler func(log *log.Logger, w http.ResponseWriter, r *http.Request) error
type App struct {
	Router           *http.ServeMux

	log              *log.Logger
	globalMiddleware []func(handler http.Handler) http.Handler
}

// New creates an App value that handle a set of routes for the application.
func New(mux *http.ServeMux, log *log.Logger, mw ...func(handler http.Handler) http.Handler) *App {
	return &App{
		Router:           mux,
		log:              log,
		globalMiddleware: mw,
	}
}

func allowGET(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
		} else {
			Respond(w, nil, http.StatusMethodNotAllowed)
		}
	})
}

func allowPOST(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		} else {
			Respond(w, nil, http.StatusMethodNotAllowed)
		}
	})
}

func allowOPTIONS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			Respond(w, nil, http.StatusMethodNotAllowed)
		}
	})
}

func CORS() func(next http.Handler) http.Handler{

}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			Respond(w, nil, http.StatusMethodNotAllowed)
		}
	})
}



// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, finalHandler http.HandlerFunc, middlwares ...func(handler http.Handler) http.Handler) {

	var h http.Handler
	h = http.HandlerFunc(finalHandler)


	// Route middle wares
	for i := len(middlwares) - 1; i >= 0; i-- {
		if middlwares[i] != nil {
			h = middlwares[i](h)
		}
	}

	if verb == http.MethodGet {
		a.globalMiddleware = append(a.globalMiddleware, allowGET)
	} else if verb == http.MethodPost {
		a.globalMiddleware = append(a.globalMiddleware, allowPOST)
	} else {
		log.Fatalln("Invalid Method")
	}


	// global middlewares
	for i := len(a.globalMiddleware) - 1; i >= 0; i-- {
		if a.globalMiddleware[i] != nil {
			h = a.globalMiddleware[i](h)
		}
	}

	a.log.Println("Adding route for", path)
	a.Router.Handle(path, h)

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

//func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	a.ServeHTTP(w, r)
//}
