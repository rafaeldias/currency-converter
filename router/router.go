package router

import (
	"fmt"
	"net/http"
	"regexp"
)

var paramsPattern = regexp.MustCompile(":([a-z]+)")

// Handler is the function used to handle http requests
type Handler func(HTTPContexter)

// Route groups the relation of url pattern and its handlers
type Route struct {
	Handlers []Handler
	Method   string
	Params   []string
	Pattern  *regexp.Regexp
}

// Router groups the router handler of requests by method (GET, POST, PUT...)
type Router struct {
	Routes      []Route
	Middlewares []Handler
}

// New returns a new Router
func New() *Router {
	return &Router{}
}

// Options parses path and set an entry in Routes with its handlers
// Same could be made for the other HTTP verbs, but for now this all we need
func (r *Router) Options(path string, handlers ...Handler) {
	r.addToRoutes(http.MethodOptions, path, handlers...)
}

// Get parses path and set an entry in Routes with its handlers
// Same could be made for the other HTTP verbs, but for now this all we need
func (r *Router) Get(path string, handlers ...Handler) {
	r.addToRoutes(http.MethodGet, path, handlers...)
}

// Use sets a new middleware
func (r *Router) Use(h Handler) {
	r.Middlewares = append(r.Middlewares, h)
}

// ServeHTTP implements the http.Handler interface in order to intercept all http calls
func (r *Router) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	for _, route := range r.Routes {
		matches := route.Pattern.FindStringSubmatch(req.URL.Path)
		if len(matches) > 0 && route.Method == req.Method {
			ctx := newHTTPContext(rw, req)

			setHTTPParams(ctx, route, matches)

			r.execMiddlewares(ctx)

			execHandlers(ctx, route)
			return
		}
	}
}

func (r *Router) execMiddlewares(ctx HTTPContexter) {
	for _, m := range r.Middlewares {
		m(ctx)
	}
}

func (r *Router) addToRoutes(method, path string, handlers ...Handler) {
	path = sanatizePath(path)
	r.Routes = append(r.Routes, Route{
		handlers,
		method,
		paramsFromPath(path),
		patternFromPath(path),
	})
}

func sanatizePath(path string) string {
	if len(path) == 0 {
		path = "/"
	}

	// Must start with a slash
	if path[0] != '/' {
		path = fmt.Sprintf("/%s", path)
	}

	return path
}

func paramsFromPath(path string) []string {
	var params []string

	if p := paramsPattern.FindAllString(path, -1); p != nil {
		for _, v := range p {
			params = append(params, v[1:])
		}
	}

	return params
}

func patternFromPath(path string) *regexp.Regexp {
	// Force exact match
	path = fmt.Sprintf("^%s$", path)
	return regexp.MustCompile(paramsPattern.ReplaceAllString(path, "([a-zA-Z0-9]+)"))
}

func setHTTPParams(ctx HTTPContexter, r Route, matches []string) {
	if len(matches) > 1 {
		for i, p := range matches[1:] {
			ctx.SetParam(r.Params[i], p)
		}
	}
}

func execHandlers(ctx HTTPContexter, r Route) {
	for _, h := range r.Handlers {
		h(ctx)
	}
}
