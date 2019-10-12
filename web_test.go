package api

import (
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteSimple(t *testing.T) {

	var app = New(nil)
	app.SetLoggingLevel(logrus.DebugLevel)

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, r, http.StatusOK, map[string]interface{}{"msg": "payload"})
	}

	app.Handle(http.MethodGet, "/test", testHandler)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
	json, _ := ioutil.ReadAll(w.Body)
	assert.JSONEq(t, `{"msg":"payload"}`, string(json))
}

func TestMiddleware(t *testing.T) {


	globalMiddle := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}
	globalMiddle2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	var app = New( globalMiddle, globalMiddle2)

	testHandler := func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, r, http.StatusOK, "hi")
	}

	testMiddlware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	testMiddlware2 := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		})
	}

	app.Handle(http.MethodGet, "/test", testHandler, testMiddlware, testMiddlware2)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	app.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
	//_, _ := ioutil.ReadAll(w.Body)
}
