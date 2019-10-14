package main

import (
	"context"
	"github.com/dbubel/api"
	"github.com/julienschmidt/httprouter"
	"time"

	"net/http"
)

type Foo struct {
	ID    string `json:"id"  validate:"required"`
	Index int    `json:"index" validate:"required"`
}

func postit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var f Foo
	err := api.UnmarshalJSON(r.Body, &f)
	if err != nil {
		api.RespondError(w, r, err, http.StatusOK)
		return
	}
	api.RespondJSON(w, r, http.StatusOK, map[string]interface{}{"status": "OK"})
}

func getit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	api.RespondJSON(w, r, http.StatusOK, &Foo{
		ID:    "123456",
		Index: 1337,
	})
}

func middlethis(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := context.WithValue(r.Context(), "middlethis", "1")
		next(w, r.WithContext(ctx), params)
	}
}

func middlethat(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := context.WithValue(r.Context(), "middlethat", "2")
		next(w, r.WithContext(ctx), params)
	}
}

func globalmiddle(next api.Handler) api.Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		next(w, r, params)
	}
}

func main() {
	apiHandlers := api.New()
	apiHandlers.GlobalMiddleware(globalmiddle)

	endpoints := api.Endpoints{
		api.NewEnpoint(http.MethodGet, "/test", getit),
		api.NewEnpoint(http.MethodGet, "/test/:paramOne", getit),
	}

	moreEndpoints := api.Endpoints{
		api.NewEnpoint(http.MethodPost, "/test", postit),
	}

	endpoints.Use(middlethis)
	moreEndpoints.Use(middlethat)
	apiHandlers.Endpoints(endpoints, moreEndpoints)

	api.StartAPI(&http.Server{
		Addr:           ":8000",
		Handler:        apiHandlers.Router,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	})
}
