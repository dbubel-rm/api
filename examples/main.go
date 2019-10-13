package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/dbubel/api"
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"time"

	"net/http"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.DebugLevel)

	app := api.NewBasic(log)



	app.GlobalMiddleware(globalmiddle)
	app.Router.RedirectTrailingSlash = true

	endpoints := api.Endpoints{
		api.NewEnpoint(http.MethodGet, "/test2", handleit),
		api.NewEnpoint(http.MethodGet, "/test2/:paramOne", handleit),
	}

	endpoints.Use(middlethat)
	endpoints.Use(middlethis)
	app.Endpoints(endpoints)


	app.StartAPI(":8000")
}

func handleit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	time.Sleep(time.Second * 8)
	fmt.Println(params.ByName("paramOne"))
	type Foo struct {
		ID    string `json:"_id"`
		Index int    `json:"index"`
	}
	var f Foo
	gofakeit.Struct(&f)
	fmt.Println("Handled")
	api.RespondJSON(w, r, http.StatusOK, f)
}

func middlethis(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		time.Sleep(time.Millisecond)
		fmt.Println("in middlethis")
		ctx := context.WithValue(r.Context(), "middlethis", "1")
		next(w, r.WithContext(ctx), params)
	}
}

func middlethat(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		time.Sleep(time.Millisecond)
		fmt.Println("in middlethat")
		ctx := context.WithValue(r.Context(), "middlethat", "2")
		next(w, r.WithContext(ctx), params)
	}
}

func globalmiddle(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		fmt.Println("in global")
		next(w, r, params)
	}
}
