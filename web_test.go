package api

import (
	"context"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestApp_SimpleRoute(t *testing.T) {
	var app = NewDefault()

	testHandler := func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		RespondJSON(w, r, http.StatusOK, map[string]interface{}{"msg": "payload"})
	}

	app.Handle(http.MethodGet, "/test", testHandler)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	resp, _ := ioutil.ReadAll(w.Body)
	assert.JSONEq(t, `{"msg":"payload"}`, string(resp))
}

func TestApp_GlobalMiddleware(t *testing.T) {
	var app = NewDefault()
	testHandler := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		RespondJSON(w, r, http.StatusOK, map[string]interface{}{"message": r.Context().Value("shared")})
	}
	app.GlobalMiddleware(middlwareOne, middlewareTwo)
	app.Handle(http.MethodGet, "/test", testHandler)
	app.Handle(http.MethodGet, "/testtwo", testHandler)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	r1 := httptest.NewRequest(http.MethodGet, "/testtwo", nil)
	w1 := httptest.NewRecorder()
	app.Router.ServeHTTP(w, r)
	app.Router.ServeHTTP(w1, r1)

	assert.Equal(t, http.StatusOK, w.Code)
	resp, _ := ioutil.ReadAll(w.Body)
	assert.JSONEq(t, `{"message":"valueonevaluetwo"}`, string(resp))

	assert.Equal(t, http.StatusOK, w1.Code)
	resp2, _ := ioutil.ReadAll(w1.Body)
	assert.JSONEq(t, `{"message":"valueonevaluetwo"}`, string(resp2))
}

func TestApp_RouteMiddleware(t *testing.T) {
	var app = NewDefault()
	testHandler := func(w http.ResponseWriter, r *http.Request, param httprouter.Params) {
		RespondJSON(w, r, http.StatusOK, map[string]interface{}{"message": r.Context().Value("shared")})
	}
	app.Handle(http.MethodGet, "/test", testHandler, middlwareOne)

	r := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	app.Router.ServeHTTP(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	resp, _ := ioutil.ReadAll(w.Body)
	assert.JSONEq(t, `{"message":"valueone"}`, string(resp))
}

func middlwareOne(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		ctx := context.WithValue(r.Context(), "globalone", "valueone")
		ctx = context.WithValue(ctx, "shared", "valueone")
		next(w, r.WithContext(ctx), params)
	}
}
func middlewareTwo(next Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		sharedValue := r.Context().Value("shared").(string)
		sharedValue = sharedValue + "valuetwo"
		ctx := context.WithValue(r.Context(), "globaltwo", "valuetwo")
		ctx = context.WithValue(ctx, "shared", sharedValue)
		next(w, r.WithContext(ctx), params)
	}
}
