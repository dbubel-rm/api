package api

import (
	"fmt"
	"io/ioutil"

	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	Success = "\u2713"
	Failed  = "\u2717"
)

func TestRouteSimple(t *testing.T) {
	mux := http.NewServeMux()

	globalMiddle := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing globalmiddle")
			next.ServeHTTP(w, r)
		})
	}

	var app = New(mux, globalMiddle)

	//testHandler := func(next http.Handler) http.Handler {
	//	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		log.Println("Executing handler")
	//		Respond(next, w, r, "hi", http.StatusOK)
	//	})
	//}

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing handler")
		Respond(w, r, "hi", http.StatusOK)
	}

	testMiddlware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing middlewareOne")
			next.ServeHTTP(w, r)
		})
	}

	testMiddlware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("Executing middlewareTwo")
			next.ServeHTTP(w, r)
		})
	}

	app.Handle(http.MethodGet, "/test", testHandler, testMiddlware, testMiddlware2)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	app.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
	b, _ := ioutil.ReadAll(w.Body)

	fmt.Println("BODYd", string(b))

}
