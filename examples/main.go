package main

import (
	"context"
	"errors"
	"github.com/dbubel/api"
	"github.com/julienschmidt/httprouter"
	"time"

	"net/http"
)

func main() {


	app := api.New()
	app.GlobalMiddleware(globalmiddle)

	endpoints := api.Endpoints{
		api.NewEnpoint(http.MethodGet, "/test", handleit),
		api.NewEnpoint(http.MethodGet, "/test/:paramOne", handleit),
	}

	endpoints.Use(middlethis)
	endpoints.Use(middlethat)
	app.Endpoints(endpoints)

	api.StartAPI(&http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	})
}

func handleit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	type Foo struct {
		ID    string `json:"_id"`
		Index int    `json:"index"`
	}

	f := Foo{
		ID:    "123456",
		Index: 1337,
	}
	_=f
	api.RespondError(w, r,errors.New("ERRO"), http.StatusBadRequest)
}

func middlethis(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		time.Sleep(time.Millisecond)
		ctx := context.WithValue(r.Context(), "middlethis", "1")
		next(w, r.WithContext(ctx), params)
	}
}

func middlethat(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		time.Sleep(time.Millisecond)
		ctx := context.WithValue(r.Context(), "middlethat", "2")
		next(w, r.WithContext(ctx), params)
	}
}

func globalmiddle(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		next(w, r, params)
	}
}
