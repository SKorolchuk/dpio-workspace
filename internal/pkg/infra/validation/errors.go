package validation

import (
	"encoding/json"
	"errors"
)

// ErrorInvalidIdentifier is used for errors when specified entity identifier format is incorrect.
var ErrorInvalidIdentifier = errors.New("identifier format is incorrect")

// ResponseError model is used to notify API consumer about errors during execution of business logic.
type ResponseError struct {
	Error            string `json:"error"`
	FieldsValidation string `json:"fields,omitempty"`
}

// RequestError model is used to notify API consumer about errors during request body validation.
type RequestError struct {
	CustomError error
	Status      int32
	Fields      error
}

// NewRequestError construct new RequestError object
func NewRequestError(e error, statusId int32) error {
	return &RequestError{
		CustomError: e,
		Status:      statusId,
		Fields:      nil,
	}
}

func (err *RequestError) Error() string {
	return err.CustomError.Error()
}

// FieldError represents error message for a specific field.
type FieldError struct {
	FieldName    string `json:"fieldName"`
	ErrorMessage string `json:"errorMessage"`
}

// FieldErrors represents a list of errors.
type FieldErrors []FieldError

// Error interface implementation for FieldErrors type.
func (fieldErrs FieldErrors) Error() string {
	fieldErrsJson, err := json.Marshal(fieldErrs)
	if err != nil {
		return err.Error()
	}

	return string(fieldErrsJson)
}

// Cause search root error from wrap chain.
func Cause(err error) error {
	cause := err
	if err == nil {
		return nil
	}

	for {
		if err = errors.Unwrap(err); err == nil {
			return cause
		}
		cause = err
	}
}
