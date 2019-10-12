package main

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	api "github.com/dbubel/api"
	"time"

	"net/http"
)

func main() {
	app := api.New()
	app.GlobalMiddleware(globalmiddle)

	endpoints := api.Endpoints{
		api.NewEnpoint(http.MethodGet, "/test2", handleit),
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

func handleit(w http.ResponseWriter, r *http.Request) {

	type Foo struct {
		ID       string `json:"_id"`
		Index    int    `json:"index"`
	}
	var f Foo
	gofakeit.Struct(&f)
	fmt.Println("Handled")
	api.RespondJSON(w, r, http.StatusOK, f)
}

func middlethis(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
		fmt.Println("in middlethis")
		ctx := context.WithValue(r.Context(), "middlethis", "1")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func middlethat(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
		fmt.Println("in middlethat")
		ctx := context.WithValue(r.Context(), "middlethat", "2")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func globalmiddle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("in global")
		next.ServeHTTP(w, r)
	})
}
