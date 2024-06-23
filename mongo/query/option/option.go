package option

import "errors"

type Option[T any] interface {
	IsSome() bool
	IsNone() bool
	Some(func(T))
	None(func())
	Match(func(T), func())
	Unwrap() (T, error)
	UnwrapOr(T) T
}

type some[T any] struct {
	value T
}

func (s some[T]) IsSome() bool {
	return true
}

func (s some[T]) IsNone() bool {
	return false
}

func (s some[T]) Some(f func(T)) {
	f(s.value)
}

func (s some[T]) None(_ func()) {
}

func (s some[T]) Match(f func(T), _ func()) {
	f(s.value)
}

func (s some[T]) Unwrap() (T, error) {
	return s.value, nil
}

func (s some[T]) UnwrapOr(t T) T {
	return s.value
}

func Some[T any](value T) Option[T] {
	return &some[T]{value: value}
}

type none[T any] struct {
}

func (n none[T]) IsSome() bool {
	return false
}

func (n none[T]) IsNone() bool {
	return true
}

func (n none[T]) Some(_ func(T)) {
}

func (n none[T]) None(f func()) {
	f()
}

func (n none[T]) Match(_ func(T), f func()) {
	f()
}

func (n none[T]) Unwrap() (T, error) {
	return *(new(T)), errors.New("option is none")
}

func (n none[T]) UnwrapOr(t T) T {
	return t
}

func None[T any]() Option[T] {
	return &none[T]{}
}
