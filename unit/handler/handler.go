package handler

import (
	"encoding/json"
	"errors"
)

type Context struct {
	byts []byte
}

func (c *Context) Bind(v interface{}) error {
	return json.Unmarshal(c.byts, v)
}

type HandlerFunc func(*Context) interface{}

type Handler struct {
	routes     map[string]HandlerFunc
	notFoundFn HandlerFunc
}

func (h *Handler) Handle(way string, fn HandlerFunc) {
	if h.routes == nil {
		h.routes = map[string]HandlerFunc{}
	}
	h.routes[way] = fn
}

func (h *Handler) NotFound(fn HandlerFunc) {
	h.notFoundFn = fn
}

func (h *Handler) Handler(way string, params []byte) (interface{}, error) {
	fun := h.notFoundFn
	if fn, ok := h.routes[way]; ok {
		fun = fn
	}

	if fun != nil {
		ret := fun(&Context{byts: params})
		if err, ok := ret.(error); ok {
			return nil, err
		}
		return ret, nil
	}
	return nil, errors.New("not found way")
}
