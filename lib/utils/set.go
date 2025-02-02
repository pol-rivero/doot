package utils

type Set[T comparable] struct {
	elements map[T]struct{}
}

func NewSet[T comparable](capacity int) Set[T] {
	return Set[T]{make(map[T]struct{}, capacity)}
}

func NewSetFromSlice[T comparable](slice []T) Set[T] {
	set := NewSet[T](len(slice))
	for _, element := range slice {
		set.Add(element)
	}
	return set
}

func (s *Set[T]) Add(element T) {
	s.elements[element] = struct{}{}
}

func (s *Set[T]) Contains(element T) bool {
	_, exists := s.elements[element]
	return exists
}

func (s *Set[T]) Remove(element T) {
	delete(s.elements, element)
}
