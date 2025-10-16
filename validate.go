package kboot

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/guestin/mob/mvalidate"
)

type Validator interface {
	Validate(i interface{}) error
	Raw() *validator.Validate
}

var _mValidator Validator

func MValidator() Validator {
	return _mValidator
}

func init() {
	var err error
	v, err := mvalidate.NewValidator(mvalidate.DefaultTranslator)
	if err != nil {
		panic(err)
	}
	_mValidator = &_validator{v: v}
}

type _validator struct {
	v mvalidate.Validator
}

func (this *_validator) Raw() *validator.Validate {
	return this.v.Raw()
}

func (this *_validator) Validate(i interface{}) error {
	k := reflect.TypeOf(i).Kind()
	if k == reflect.Struct || (k == reflect.Ptr && reflect.ValueOf(i).Elem().Kind() == reflect.Struct) {
		return this.v.Struct(i)
	}
	if k == reflect.Slice {
		return this.v.Var(i, "required,dive,required")
	}
	return this.v.Var(i, "required")
}
