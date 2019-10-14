package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)
var log100s = true
var log200s = true
var log300s = true
var log400s = true
var log500s = true
var apiLogger *logrus.Logger

func Log100s(val bool) {
	log100s = val
}
func Log200s(val bool) {
	log200s = val
}
func Log300s(val bool) {
	log300s = val
}
func Log400s(val bool) {
	log400s = val
}
func Log500s(val bool) {
	log500s = val
}

func logHandler(r *http.Request, code int) {
	printLog := func() {
		apiLogger.WithFields(logrus.Fields{
			"method":     r.Method,
			"url":        r.RequestURI,
			"contentLen": r.ContentLength,
			"ms":         time.Now().Sub(r.Context().Value("ts").(time.Time)).Milliseconds(),
			"code":       code,
		}).Info()
	}

	if code >= 100 && code < 200 && log100s {
		printLog()
	} else if code >= 200 && code < 300 && log200s {
		printLog()
	} else if code >= 300 && code < 400 && log300s {
		printLog()
	} else if code >= 400 && code < 500 && log400s {
		printLog()
	} else if code >= 500 && code < 500 && log500s {
		printLog()
	}
}