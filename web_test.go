package api

import (
	"github.com/julienschmidt/httprouter"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRouteSimple(t *testing.T) {

	var app = NewBasic()
	app.SetLoggingLevel(log.InfoLevel)

	testHandler := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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

	globalMiddle := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			next(w, r, params)
		}
	}
	globalMiddle2 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			next(w, r, params)
		}
	}

	var app = NewBasic()
	app.GlobalMiddleware(globalMiddle, globalMiddle2)

	testHandler := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		RespondJSON(w, r, http.StatusOK, "hi")
	}

	testMiddlware := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			next(w, r, params)
		}
	}

	testMiddlware2 := func(next Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
			next(w, r, params)
		}
	}

	app.Handle(http.MethodGet, "/test", testHandler, testMiddlware, testMiddlware2)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	app.Router.ServeHTTP(w, r)
	assert.Equal(t, http.StatusOK, w.Code, "Response code should be 200")
	//_, _ := ioutil.ReadAll(w.Body)
}
