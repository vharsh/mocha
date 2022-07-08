package to

// Repeat returns true if total request hits for current mock is equal or lower total the provided max call times.
// If Repeat is used direct, it must be set using Mock After Expectations.
func Repeat(times int) Matcher[any] {
	count := 0

	m := Matcher[any]{}
	m.Name = "Repeat"
	m.Matches = func(_ any, params Args) (bool, error) {
		count++

		return count <= times, nil
	}

	return m
}