package util

type Set[T comparable] struct {
	elems map[T]struct{}
}

func NewSet[T comparable](keys ...T) Set[T] {
	s := Set[T]{elems: make(map[T]struct{})}
	return s
}

func (s Set[T]) Add(key T) {
	s.elems[key] = struct{}{}
}

func (s Set[T]) Remove(key T) {
	delete(s.elems, key)
}

func (s Set[T]) Contains(key T) bool {
	_, ok := s.elems[key]
	return ok
}

func (s Set[T]) Len() int {
	return len(s.elems)
}
