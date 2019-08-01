package api

import (
	"encoding/json"
	"net/http"
)

func RespondError(next http.Handler, w http.ResponseWriter, r *http.Request, err error, code int) {
	Respond(next, w, r, rr{"error": err.Error()}, code)
}

func Respond(next http.Handler, w http.ResponseWriter, r *http.Request, data interface{}, code int) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		RespondError(next, w, r, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(jsonData)
	next.ServeHTTP(w, r)
}
