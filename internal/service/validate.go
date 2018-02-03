package service

import (
	"encoding/json"
	"net/http"

	valid "github.com/asaskevich/govalidator"
)

// decodeAndValidae performs JSON decoding from an HTTP request and validates it using govalidator annotations
func decodeAndValidate(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}

	_, err := valid.ValidateStruct(v)
	return err
}
