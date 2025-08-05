package utils

import (
	"reflect"
	"strings"
)

func TrimStrings(s any) {
	v := reflect.ValueOf(s).Elem()

	for i := range v.NumField() {
		field := v.Field(i)
		if field.Kind() == reflect.String {
			field.SetString(strings.TrimSpace(field.String()))
		}
	}
}
