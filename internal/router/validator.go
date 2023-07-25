package router

import "gopkg.in/go-playground/validator.v9"

// NewValidator - creating new validator
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

//Validator - custom validator struct
type Validator struct {
	validator *validator.Validate
}

// Validate method for calling validation
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
