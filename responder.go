package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

var Log100s = true
var Log200s = true
var Log300s = true
var Log400s = true
var Log500s = true

func RespondError(w http.ResponseWriter, r *http.Request, err error, code int, description ...string) {
	switch err.(type) {
	default:
		RespondJSON(w, r, code, map[string]interface{}{"error": err.Error(), "description": description})
	case InvalidError:
		RespondJSON(w, r, code, err)
	}
}

func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		jsonData, err := json.Marshal(data)
		if err != nil {
			RespondError(w, r, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(jsonData)
	}
	LogRequest(apiLogger, r, code)
}

func Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	} else {
		contentType := http.DetectContentType(data)
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
		w.Write(data)
	}
	LogRequest(apiLogger, r, code)
}

func LogRequest(log *logrus.Logger, r *http.Request, code int, description ...string) {
	printLog := func() {
		log.WithFields(logrus.Fields{
			"method":     r.Method,
			"url":        r.RequestURI,
			"contentLen": r.ContentLength,
			"ms":         time.Now().Sub(r.Context().Value("ts").(time.Time)).Milliseconds(),
			"code":       code,
		}).Info(description)
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
