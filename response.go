package api

import (
	"encoding/json"
	"net/http"
)

func (a *App) RespondError(w http.ResponseWriter, r *http.Request, err error, code int, description ...string) {
	switch err.(type) {
	default:
		a.RespondJSON(w, r, code, map[string]interface{}{"error": err.Error(), "description": description})
	case InvalidError:
		a.RespondJSON(w, r, code, err)
	}
}

func (a *App) RespondJSON(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	} else {
		jsonData, err := json.Marshal(data)
		if err != nil {
			a.RespondError(w, r, err, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(jsonData)
	}
	a.logRequest(r, code)
}

func (a *App) Respond(w http.ResponseWriter, r *http.Request, code int, data []byte) {
	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	} else {
		contentType := http.DetectContentType(data)
		w.Header().Set("Content-Type", contentType)
		w.WriteHeader(code)
		w.Write(data)
	}
	a.logRequest(r, code)
}
