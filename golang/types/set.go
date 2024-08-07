package cav

import "fmt"

type ISet[T comparable] interface {
	Add(value T)
	Contains(value T) bool
	ForEach(f func(T))
	Copy() ISet[T]
	Union(other ISet[T]) ISet[T]
	Intersect(other ISet[T]) ISet[T]
	Minus(other ISet[T]) ISet[T]
	Equals(other ISet[T]) bool
	String() string
}

type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(value T) {
	s[value] = struct{}{}
}

func (s Set[T]) Contains(value T) bool {
	_, ok := s[value]
	return ok
}

func (s Set[T]) ForEach(f func(T)) {
	for value := range s {
		f(value)
	}
}

func (s Set[T]) Copy() ISet[T] {
	result := MakeSet[T]()
	s.ForEach(func(value T) {
		result.Add(value)
	})
	return result
}

func (s Set[T]) Union(other ISet[T]) ISet[T] {
	result := s.Copy()
	other.ForEach(func(value T) {
		result.Add(value)
	})
	return result
}

func (s Set[T]) Intersect(other ISet[T]) ISet[T] {
	result := MakeSet[T]()
	s.ForEach(func(value T) {
		if other.Contains(value) {
			result.Add(value)
		}
	})
	return result
}

func (s Set[T]) Minus(other ISet[T]) ISet[T] {
	result := MakeSet[T]()
	s.ForEach(func(value T) {
		if !other.Contains(value) {
			result.Add(value)
		}
	})
	return result
}

func (s Set[T]) Equals(other ISet[T]) bool {
	if other == nil {
		return false
	}
	result := true
	s.ForEach(func(value T) {
		if !other.Contains(value) {
			result = false
		}
	})
	other.ForEach(func(value T) {
		if !s.Contains(value) {
			result = false
		}
	})
	return result
}

func (s Set[T]) String() string {
	result := "{"
	first := true
	s.ForEach(func(value T) {
		if !first {
			result += ", "
		} else {
			first = false
		}
		result += fmt.Sprint(value)
	})
	result += "}"
	return result
}

func MakeSet[T comparable]() ISet[T] {
	return Set[T]{}
}

func MakeSetOf[T comparable](values ...T) ISet[T] {
	set := MakeSet[T]()
	for _, v := range values {
		set.Add(v)
	}
	return set
}
