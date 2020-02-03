package persist

import (
	"reflect"
	"sync"
)

var (
	initOnce        sync.Once
	defaultRegistry *Registry
)

func DefaultRegistry() *Registry {
	initOnce.Do(func() {
		defaultRegistry = NewRegistry()
	})
	return defaultRegistry
}

type Registry struct {
	sync.RWMutex
	reg map[reflect.Type]interface{}
}

func NewRegistry() *Registry {
	return &Registry{
		sync.RWMutex{},
		make(map[reflect.Type]interface{}),
	}
}

func (r *Registry) Set(t reflect.Type, p interface{}) {
	t = indirect(t)
	r.Lock()
	defer r.Unlock()
	r.reg[t] = p
}

func (r *Registry) SetOnce(t reflect.Type, p interface{}) {
	t = indirect(t)
	r.Lock()
	defer r.Unlock()
	if _, ok := r.reg[t]; !ok {
		r.reg[t] = p
	}
}

func (r *Registry) Get(t reflect.Type) (interface{}, bool) {
	t = indirect(t)
	r.RLock()
	defer r.RUnlock()
	p, ok := r.reg[t]
	return p, ok
}

func (r *Registry) GetFor(v interface{}) (interface{}, bool) {
	return r.Get(reflect.ValueOf(v).Type())
}

func indirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}
