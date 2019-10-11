package api

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"

	"net/http"
	"time"
)

func RespondError(w http.ResponseWriter, r *http.Request, err error, code int) {
	Respond(w, r, map[string]interface{}{"error": err.Error()}, code)
}

func Respond(w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		RespondError(w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)

	responseTime := time.Now().Sub(r.Context().Value("ts").(time.Time))

	l.WithFields(log.Fields{"method": r.Method,
		"url":          r.RequestURI,
		"contentLenth": r.ContentLength,
		"ip":           r.RemoteAddr,
		"ts":           responseTime.String(),
		"code":         code,
	}).Info("request details")

	//next.ServeHTTP(w, r)
}
