package api

import (
	"encoding/json"
	"net/http"
)

// JSONError is the response for errors that occur within the API.
type JSONError struct {
	Error  string       `json:"error"`
	Fields InvalidError `json:"fields,omitempty"`
}

// RespondError sends JSON describing the error
func RespondError(next http.Handler, w http.ResponseWriter, r *http.Request, err error, code int) {
	Respond(next, w, r, JSONError{Error: err.Error()}, code)
}

// Respond sends JSON to the client.
// If code is StatusNoContent, v is expected to be nil.
func Respond(next http.Handler, w http.ResponseWriter, r *http.Request, data interface{}, code int) {

	if code == http.StatusNoContent || data == nil {
		w.WriteHeader(code)
		return
	}

	// Marshal the data into a JSON string.
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		// Should respond with internal server error.
		RespondError(next, w, r, err, http.StatusInternalServerError)
		return
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response and context.
	w.WriteHeader(code)

	// Send the result back to the client.
	w.Write(jsonData)

	next.ServeHTTP(w, r)

}
