package mapx

import "zimg/utils"

type SliceMap[K comparable, V any] struct {
	m map[K]V
	s []K
}

func NewSliceMap[K comparable, V any]() *SliceMap[K, V] {
	return &SliceMap[K, V]{
		m: make(map[K]V),
		s: make([]K, 0),
	}
}

func (r *SliceMap[K, V]) Len() int {
	return len(r.m)
}

func (r *SliceMap[K, V]) Keys() []K {
	return r.s
}

func (r *SliceMap[K, V]) Values() []V {
	values := make([]V, 0)
	for _, v := range r.s {
		values = append(values, r.m[v])
	}
	return values
}

func (r *SliceMap[K, V]) Set(key K, value V) {
	_, ok := r.m[key]
	r.m[key] = value
	if ok {
		i := utils.IndexOf(r.s, key)
		r.s[i] = key
	} else {
		r.s = append(r.s, key)
	}
}

func (r *SliceMap[K, V]) Get(key K) (V, bool) {
	value, ok := r.m[key]
	return value, ok
}

func (r *SliceMap[K, V]) Delete(key K) {
	_, ok := r.m[key]
	delete(r.m, key)
	if ok {
		i := utils.IndexOf(r.s, key)
		r.s = append(r.s[:i], r.s[i+1:]...)
	}
}
