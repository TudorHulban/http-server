package router

import "net/http"

type Route struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

type Router struct {
	routes map[string]func(http.ResponseWriter, *http.Request)
}

func NewRouter() *Router {
	return &Router{
		routes: make(map[string]func(http.ResponseWriter, *http.Request)),
	}
}

func (r *Router) AddRoute(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	r.routes[pattern] = handler
}

func (r *Router) FindHandler(path string) (func(http.ResponseWriter, *http.Request), bool) {
	handler, ok := r.routes[path]
	return handler, ok
}
