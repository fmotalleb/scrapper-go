package utils

import (
	"fmt"
)

type (
	getter = func() string
	Vars   map[string]varValue
)

type varValue struct {
	isGenerative bool
	value        string
	get          getter
}

func (v *varValue) getValue() string {
	if v.isGenerative {
		return v.get()
	} else {
		return v.value
	}
}

func (v Vars) Snapshot() map[string]string {
	snap := make(map[string]string)
	for k, g := range v {
		snap[k] = g.getValue()
	}
	return snap
}

func (v Vars) LiveSnapshot() map[string]any {
	snap := make(map[string]any)
	for k, g := range v {
		if g.isGenerative {
			snap[k] = g.get
		} else {
			snap[k] = g.value
		}
	}
	return snap
}

func (v Vars) SetOnce(key string, value string) {
	v[key] = varValue{
		isGenerative: false,
		value:        value,
	}
}

func (v Vars) SetGetter(key string, getter getter) {
	v[key] = varValue{
		isGenerative: true,
		get:          getter,
	}
}

func (v Vars) Get(key string) (string, bool) {
	item, ok := v[key]

	return item.getValue(), ok
}

func (v Vars) GetOr(key string, def string) string {
	value, ok := v.Get(key)
	if ok {
		return value
	}
	return def
}

func (v Vars) GetOrFail(key string) (string, error) {
	value, ok := v.Get(key)
	if ok {
		return value, nil
	}
	return "", fmt.Errorf("use of undefined variable: %s", key)
}
