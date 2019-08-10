package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteSimple(t *testing.T) {
	mux := httprouter.New()
	var app = New(mux, nil)
	app.SetLoggingLevel(logrus.WarnLevel)

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		Respond(w, r, map[string]interface{}{"msg":"payload"}, http.StatusOK)
	}

	app.Handle(http.MethodGet, "/test", testHandler)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
	json, _:= ioutil.ReadAll(w.Body)
	assert.JSONEq(t,`{"msg":"payload"}`,string(json) )
}

//
//func TestMiddleware(t *testing.T) {
//	mux := httprouter.New()
//
//	globalMiddle := func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			log.Println("Executing globalmiddle")
//			next.ServeHTTP(w, r)
//		})
//	}
//
//	var app = New(mux, globalMiddle)
//
//	testHandler := func(w http.ResponseWriter, r *http.Request) {
//		log.Println("Executing handler")
//		Respond(w, r, "hi", http.StatusOK)
//	}
//
//	testMiddlware := func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			log.Println("Executing middlewareOne")
//			next.ServeHTTP(w, r)
//		})
//	}
//
//	testMiddlware2 := func(next http.Handler) http.Handler {
//		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//			log.Println("Executing middlewareTwo")
//			next.ServeHTTP(w, r)
//		})
//	}
//
//	app.Handle(http.MethodGet, "/test", testHandler, testMiddlware, testMiddlware2)
//
//	r := httptest.NewRequest(http.MethodGet, "/test", nil)
//	w := httptest.NewRecorder()
//
//	app.Router.ServeHTTP(w, r)
//	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
//	b, _ := ioutil.ReadAll(w.Body)
//
//	fmt.Println("BODYd", string(b))
//
//}

