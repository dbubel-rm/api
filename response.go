package api

import (
	"encoding/json"
	"net/http"
)

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
	logRequest(r, code)
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
	logRequest(r, code)
}
