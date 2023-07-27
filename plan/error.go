package plan

import (
	"fmt"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
)

func newValidateError(tag string, field string, value interface{}, errFormat string, args ...any) validateError {
	valueOf := reflect.ValueOf(value)
	return validateError{
		tag:   tag,
		field: field,
		value: value,
		kind:  valueOf.Kind(),
		typ:   valueOf.Type(),
		err:   fmt.Sprintf(errFormat, args...),
	}
}

var _ validator.FieldError = validateError{}

type validateError struct {
	tag   string
	field string
	value interface{}
	kind  reflect.Kind
	typ   reflect.Type
	err   string
}

func (v validateError) Tag() string {
	return v.tag
}

func (v validateError) ActualTag() string {
	return v.tag
}

func (v validateError) Namespace() string {
	return ""
}

func (v validateError) StructNamespace() string {
	return ""
}

func (v validateError) Field() string {
	return v.field
}

func (v validateError) StructField() string {
	return v.field
}

func (v validateError) Value() interface{} {
	return v.value
}

func (v validateError) Param() string {
	return ""
}

func (v validateError) Kind() reflect.Kind {
	return v.kind
}

func (v validateError) Type() reflect.Type {
	return v.typ
}

func (v validateError) Translate(ut ut.Translator) string {
	return v.err
}

func (v validateError) Error() string {
	return v.err
}
