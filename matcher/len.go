package matcher

import "reflect"

func Len[E any](length int) Matcher[E] {
	return func(v E, params Params) (bool, error) {
		value := reflect.ValueOf(v)
		return value.Len() == length, nil
	}
}