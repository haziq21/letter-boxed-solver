package sets

type Set[T comparable] map[T]struct{}

func New[T comparable](items ...T) Set[T] {
	s := make(Set[T])
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

func (s Set[T]) Contains(item T) bool {
	_, ok := s[item]
	return ok
}

func (s Set[T]) Remove(item T) {
	delete(s, item)
}

func (s Set[T]) ToSlice() []T {
	result := make([]T, 0, len(s))
	for k := range s {
		result = append(result, k)
	}
	return result
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	res := New[T]()
	for k := range s {
		res[k] = struct{}{}
	}
	for k := range other {
		res[k] = struct{}{}
	}
	return res
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	res := New[T]()
	for k := range s {
		if other.Contains(k) {
			res[k] = struct{}{}
		}
	}
	return res
}

// Diff returns a new set containing elements in s that are not in other.
func (s Set[T]) Diff(other Set[T]) Set[T] {
	res := New[T]()
	for k := range s {
		if !other.Contains(k) {
			res[k] = struct{}{}
		}
	}
	return res
}

func (s Set[T]) IsSubsetOf(other Set[T]) bool {
	for k := range s {
		if !other.Contains(k) {
			return false
		}
	}
	return true
}
