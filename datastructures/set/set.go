package set

type Set[T comparable] struct {
	data map[T]struct{}
}

func New[T comparable](capacity int) Set[T] {
	return Set[T]{
		data: make(map[T]struct{}, capacity),
	}
}

func Define[T comparable](items ...T) Set[T] {
	tr := New[T](len(items))
	tr.AddFromSlice(items)
	return tr
}

func (s *Set[T]) Add(item ...T) {
	for idx := range item {
		s.data[item[idx]] = struct{}{}
	}
}

func (s *Set[T]) Contains(item T) bool {
	_, exists := s.data[item]
	return exists
}

func (s *Set[T]) Remove(item T) {
	delete(s.data, item)
}

func (s *Set[T]) Len() int {
	return len(s.data)
}

func (s *Set[T]) Range(callback func(item T) bool) {
	for item := range s.data {
		if !callback(item) {
			return
		}
	}
}

func (s *Set[T]) Clone() Set[T] {
	tr := New[T](s.Len())
	tr.AddFromAnother(*s)
	return tr
}

func (s *Set[T]) Union(another Set[T]) {
	s.AddFromAnother(another)
}

func (s *Set[T]) Intersect(another Set[T]) {
	s.Range(func(item T) bool {
		if !another.Contains(item) {
			s.Remove(item)
		}
		return true
	})
}

func (s *Set[T]) AddFromAnother(another Set[T]) {
	another.Range(func(item T) bool {
		s.Add(item)
		return true
	})
}

func (s *Set[T]) AddFromSlice(items []T) {
	for _, item := range items {
		s.Add(item)
	}
}

func (s *Set[T]) ToSlice() []T {
	tr := make([]T, 0, s.Len())
	s.Range(func(item T) bool {
		tr = append(tr, item)
		return true
	})
	return tr
}
