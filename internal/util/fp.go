package util

import "cmp"

func Map[T, K any](f func(T) K, s ...T) []K {
	ret := make([]K, 0, len(s))
	for _, v := range s {
		ret = append(ret, f(v))
	}
	return ret
}

func MapErr[T, K any](f func(T) (K, error), s ...T) (result []K, errIndex int, err error) {
	ret := make([]K, 0, len(s))
	for i, v := range s {
		r, err := f(v)
		if err != nil {
			return nil, i, err
		}
		ret = append(ret, r)
	}
	return ret, 0, nil
}

func Filter[T any](f func(T) bool, s ...T) []T {
	ret := make([]T, 0, len(s))
	for _, v := range s {
		if f(v) {
			ret = append(ret, v)
		}
	}
	return ret
}

func Keys[K cmp.Ordered, V any](m map[K]V) []K {
	ret := make([]K, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

func Values[K cmp.Ordered, V any](m map[K]V) []V {
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

type Result[T any] struct {
	Value T
	Err   error
}

func Ok[T any](v T) Result[T] {
	return Result[T]{Value: v}
}

func Err[T any](e error) Result[T] {
	return Result[T]{Err: e}
}
