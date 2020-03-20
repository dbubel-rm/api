package api

import "fmt"

// Invalid describes a validation error belonging to a specific field.
type Invalid struct {
	Fld string `json:"field_name"`
	Err string `json:"error"`
}

// InvalidError is a custom error type for invalid fields.
type InvalidError []Invalid

// Error implements the error interface for InvalidError.
func (err InvalidError) Error() string {
	var str string
	for _, v := range err {
		str = fmt.Sprintf("%s,{%s:%s}", str, v.Fld, v.Err)
	}
	return str
}
