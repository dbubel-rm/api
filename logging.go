package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var Log100s = true
var Log200s = true
var Log300s = true
var Log400s = true
var Log500s = true


func (a *App) logRequest(r *http.Request, code int) {
	printLog := func() {
		a.ApiLogger.WithFields(logrus.Fields{
			"method":     r.Method,
			"url":        r.RequestURI,
			"contentLen": r.ContentLength,
			"ms":         time.Now().Sub(r.Context().Value("ts").(time.Time)).Milliseconds(),
			"code":       code,
		}).Info()
	}

	if code >= 100 && code < 200 && Log100s {
		printLog()
	} else if code >= 200 && code < 300 && Log200s {
		printLog()
	} else if code >= 300 && code < 400 && Log300s {
		printLog()
	} else if code >= 400 && code < 500 && Log400s {
		printLog()
	} else if code >= 500 && code < 500 && Log500s {
		printLog()
	}
}
