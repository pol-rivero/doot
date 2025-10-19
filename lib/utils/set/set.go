package set

type Set[T comparable] struct {
	elements map[T]struct{}
}

func New[T comparable](capacity int) Set[T] {
	return Set[T]{make(map[T]struct{}, capacity)}
}

func NewFromSlice[T comparable](slice []T) Set[T] {
	set := New[T](len(slice))
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

func (s *Set[T]) Iter() map[T]struct{} {
	return s.elements
}

func (s *Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s.elements))
	for element := range s.elements {
		slice = append(slice, element)
	}
	return slice
}
