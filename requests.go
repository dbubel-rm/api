package api

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
)

// validate provides a validator for checking models.
var validate = validator.New()

// Unmarshal decodes the input to the struct type and checks the
// fields to verify the value is in a proper state.
func UnmarshalJSON(r io.Reader, v interface{}) error {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		return err
	}

	var inv InvalidError
	if fve := validate.Struct(v); fve != nil {
		for _, fe := range fve.(validator.ValidationErrors) {
			inv = append(inv, Invalid{Fld: fe.Field(), Err: fe.Tag()})
		}
		return inv
	}
	return nil
}
