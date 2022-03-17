package enumerable

type Enumerable[K comparable, V any] struct {
	in chan entry[K, V]
}

type entry[K comparable, V any] struct {
	key K
	val V
}

func FromSlice[TSlice ~[]T, T any](slice TSlice) Enumerable[int, T] {
	in := make(chan entry[int, T])
	go func() {
		for i, v := range slice {
			in <- entry[int, T]{key: i, val: v}
		}

		close(in)
	}()

	return Enumerable[int, T]{in}
}

func FromMap[Map ~map[K]V, K comparable, V any](m Map) Enumerable[K, V] {
	in := make(chan entry[K, V])
	go func() {
		for k, v := range m {
			in <- entry[K, V]{key: k, val: v}
		}

		close(in)
	}()

	return Enumerable[K, V]{in}
}

func (enum Enumerable[K, V]) Next() (K, V, bool) {
	next, ok := <-enum.in
	return next.key, next.val, ok
}

func (enum Enumerable[K, V]) Values() []V {
	values := make([]V, 0)
	for e := range enum.in {
		values = append(values, e.val)
	}

	return values
}

func (enum Enumerable[K, V]) Keys() []K {
	keys := make([]K, 0)
	for e := range enum.in {
		keys = append(keys, e.key)
	}

	return keys
}

func (enum Enumerable[K, V]) ToMap() map[K]V {
	m := make(map[K]V)
	for e := range enum.in {
		m[e.key] = e.val
	}

	return m
}

func (enum Enumerable[K, V]) Filter(predicate func(K, V) bool) Enumerable[K, V] {
	in := make(chan entry[K, V])
	go func() {
		for e := range enum.in {
			if predicate(e.key, e.val) {
				in <- e
			}
		}

		close(in)
	}()

	return Enumerable[K, V]{in}
}

func Map[K comparable, V any, L comparable, U any](enum Enumerable[K, V], fn func(K, V) (L, U)) Enumerable[L, U] {
	in := make(chan entry[L, U])
	go func() {
		for e := range enum.in {
			l, u := fn(e.key, e.val)
			in <- entry[L, U]{key: l, val: u}
		}

		close(in)
	}()

	return Enumerable[L, U]{in}
}

func MapValues[K comparable, V any, U any](enum Enumerable[K, V], fn func(V) U) Enumerable[K, U] {
	in := make(chan entry[K, U])
	go func() {
		for e := range enum.in {
			u := fn(e.val)
			in <- entry[K, U]{key: e.key, val: u}
		}

		close(in)
	}()

	return Enumerable[K, U]{in}
}
