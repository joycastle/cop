package util

import (
	"reflect"
)

func IsVarPointor(v interface{}) bool {
	return reflect.ValueOf(&v).Kind() == reflect.Ptr
}
