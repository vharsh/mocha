package matcher

func Not[V any](m Matcher[V]) Matcher[V] {
	return func(v V, params Params) (bool, error) {
		result, err := m(v, params)
		return !result, err
	}
}