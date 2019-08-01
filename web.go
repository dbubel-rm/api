package api

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"
)

//type Handler func(log *log.Logger, w http.ResponseWriter, r *http.Request) error
type App struct {
	Router           *http.ServeMux
	log              *log.Logger
	globalMiddleware []func(handler http.Handler) http.Handler
}

// New creates an App value that handle a set of routes for the application.
func New(mux *http.ServeMux, mw ...func(handler http.Handler) http.Handler) *App {
	return &App{
		Router:           mux,
		log:              log.New(os.Stdout, "API", log.LstdFlags|log.Lmicroseconds|log.Lshortfile),
		globalMiddleware: mw,
	}
}

func allowGET(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			next.ServeHTTP(w, r)
		} else {
			Respond(next, w,r, nil, http.StatusMethodNotAllowed)
		}
	})
}

func allowPOST(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			next.ServeHTTP(w, r)
		} else {
			Respond(next, w, r,nil, http.StatusMethodNotAllowed)
		}
	})
}

func requestStart(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//fmt.Println("Start")
		//next.ServeHTTP(w, r)
		ctx := context.WithValue(r.Context(), "ts", time.Now())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

//func requestEnd(next http.Handler) http.Handler {
//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		fmt.Println("End")
//		next.ServeHTTP(w, r)
//	})
//}


func X(w http.ResponseWriter, r *http.Request) {
	log.Printf("[API][%s] %v", r.Method, time.Now().Sub(r.Context().Value("ts").(time.Time)))
}

func allowOPTIONS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
		} else {
			Respond(next,w, r,nil, http.StatusMethodNotAllowed)
		}
	})
}

// Handle is our mechanism for mounting Handlers for a given HTTP verb and path
// pair, this makes for really easy, convenient routing.
func (a *App) Handle(verb, path string, finalHandler func(handler http.Handler) http.Handler, middlwares ...func(handler http.Handler) http.Handler) {

	//testHandler := func(w http.ResponseWriter, r *http.Request) {
	//	log.Println("Executing handler")
	//	Respond(w, "hi", http.StatusOK)
	//}

	//testMiddlware2 := func(next http.Handler) http.Handler {
	//	return finalHandler
	//}

	var h http.Handler
	h = http.HandlerFunc(X)

	// Route middle wares
	middlwares = append(middlwares, finalHandler)

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

	a.globalMiddleware = append([]func(handler http.Handler) http.Handler{requestStart}, a.globalMiddleware...)
	// global middlewares
	for i := len(a.globalMiddleware) - 1; i >= 0; i-- {
		if a.globalMiddleware[i] != nil {
			h = a.globalMiddleware[i](h)
		}
	}
	//middlwares = append(middlwares, requestStart)

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

//func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	a.ServeHTTP(w, r)
//}
