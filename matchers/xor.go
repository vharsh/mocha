package matchers

// XOR is a exclusive or matcher
func XOR[V any](first Matcher[V], second Matcher[V]) Matcher[V] {
	m := Matcher[V]{}
	m.Name = "Xor"
	m.Matches = func(v V, args Args) (bool, error) {
		a, err := first.Matches(v, args)
		if err != nil {
			return false, err
		}

		b, err := second.Matches(v, args)
		if err != nil {
			return false, err
		}

		return a != b, nil
	}

	return m
}
