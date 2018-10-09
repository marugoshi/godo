package Router

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
)

const (
	TextPlain = "text/plain"
	TextHtml = "text/html"
)

type Handler func(e *Exchange)

type Route struct {
	Pattern string
	Method string
	Handler
	Re *regexp.Regexp
}

type Router struct {
	Routes       []Route
	ErrorHandler Handler
}

func New(contentType string) *Router {
	var errorHandler Handler
	switch contentType {
	case TextPlain:
		errorHandler = func(e *Exchange) {
			e.Text(http.StatusNotFound, "Not Found")
		}
	case TextHtml:
		errorHandler = func(e *Exchange) {
		}
	default:
		errorHandler = func(e *Exchange) {
		}
	}

	router := &Router{}
	router.ErrorHandler = errorHandler

	for _, route := range Routes {
		route.Re = regexp.MustCompile(route.Pattern)
		router.Routes = append(router.Routes, route)
	}

	return router
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	e := &Exchange{ResponseWriter: res, Request: req}
	for _, route := range r.Routes {
		if matches := route.Re.FindStringSubmatch(e.URL.Path); len(matches) > 0 {
			if len(matches) > 1 {
				e.QueryString = matches[1:]
			}
			route.Handler(e)
			return
		}
	}
	r.ErrorHandler(e)
}

type Exchange struct {
	http.ResponseWriter
	*http.Request
	QueryString []string
}

func (e *Exchange) Text(code int, body string) {
	e.ResponseWriter.Header().Set("Content-Type", "text/plain")
	e.WriteHeader(code)
	io.WriteString(e.ResponseWriter, fmt.Sprintf("%s\n", body))
}