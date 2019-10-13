package api

import (
	"encoding/json"
	"github.com/sirupsen/logrus"

	"net/http"
	"time"
)

func RespondError(w http.ResponseWriter, r *http.Request, err error, code int, description ...string) {
	RespondJSON(w, r, code, map[string]interface{}{"error": err.Error(), "description": description})
}
func logHandler(r *http.Request, code int) {
	apiLogger.WithFields(logrus.Fields{
		"method":     r.Method,
		"url":        r.RequestURI,
		"contentLen": r.ContentLength,
		"ip":         r.RemoteAddr,
		"ms":         time.Now().Sub(r.Context().Value("ts").(time.Time)).Milliseconds(),
		"code":       code,
	}).Info()
}
func RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		jsonData, err := json.MarshalIndent(data, "", "  ")
		if err != nil {
			RespondError(w, r, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(jsonData)
	}
	logHandler(r, code)
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
	logHandler(r, code)
}
