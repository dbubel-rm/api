package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	api "github.com/dbubel/api"
	"github.com/julienschmidt/httprouter"
	"time"

	"net/http"
)

func main() {
	app := api.NewBasic()
	app.GlobalMiddleware(globalmiddle)

	endpoints := api.Endpoints{
		api.NewEnpoint(http.MethodGet, "/test2", handleit),
		api.NewEnpoint(http.MethodGet, "/test2/:paramOne", handleit),
	}

	endpoints.Use(middlethat)
	endpoints.Use(middlethis)
	app.Endpoints(endpoints)

	a := http.Server{
		Addr:           ":8000",
		Handler:        app.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Println("running...")
	a.ListenAndServe()
}

func handleit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

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
