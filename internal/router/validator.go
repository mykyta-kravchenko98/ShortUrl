package router

import "gopkg.in/go-playground/validator.v9"

//Creating new validator
func NewValidator() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}

//Custom validator struct
type Validator struct {
	validator *validator.Validate
}

// Method for calling validation
func (v *Validator) Validate(i interface{}) error {
	return v.validator.Struct(i)
}
