package slices

// RemoveInPlace removes all elements from a slice that match the given predicate.
// Does not allocate a new slice.
func RemoveInPlace[T any](collection []T, predicate func(T, int) bool) []T {
	i := 0
	var t T
	for j, x := range collection {
		if !predicate(x, j) {
			collection[j] = t
			collection[i] = x
			i++
		}
	}
	return collection[:i]
}
