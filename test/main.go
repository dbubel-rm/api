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

	endpoints := api.Endpoints{
		api.NewEnpoint(http.MethodGet, "/test2", handleit),
	}

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
		GUID     string `json:"guid"`
		IsActive bool   `json:"isActive"`
		Balance  string `json:"balance"`
		Picture  string `json:"picture"`
		Age      int    `json:"age"`
		EyeColor string `json:"eyeColor"`
		Name     struct {
			First string `json:"first"`
			Last  string `json:"last"`
		} `json:"name"`
		Company    string   `json:"company"`
		Email      string   `json:"email"`
		Phone      string   `json:"phone"`
		Address    string   `json:"address"`
		About      string   `json:"about"`
		Registered string   `json:"registered"`
		Latitude   string   `json:"latitude"`
		Longitude  string   `json:"longitude"`
		Tags       []string `json:"tags"`
		Range      []int    `json:"range"`
		Friends    []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"friends"`
		Greeting      string `json:"greeting"`
		FavoriteFruit string `json:"favoriteFruit"`
	}

	// Pass your struct as a pointer
	var f Foo
	gofakeit.Struct(&f)

	api.RespondJSON(w, r, http.StatusOK, nil)
}

func middlethis(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(time.Millisecond)
		fmt.Println("in middle")
		ctx := context.WithValue(r.Context(), "something", "value")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
