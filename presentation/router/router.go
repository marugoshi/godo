package router

import (
	"github.com/marugoshi/gobm/presentation/httputils"
	"net/http"
	"regexp"
)

type Routes interface {
	Data() []Route
}

type routes struct {
	data []Route
}

func (r *routes) Data() []Route {
	return r.data
}

func NewRoutes() Routes {
	routes := &routes{}
	return routes
}

type Route struct {
	Pattern *regexp.Regexp
	Method  string
	httputils.Handler
}

type Router struct {
	Routes
	httputils.Handler
}

func NewRouter(contentType string) *Router {
	return &Router{NewRoutes(), notFoundErrorHandler(contentType)}
}

func notFoundErrorHandler(contentType string) httputils.Handler {
	var errorHandler httputils.Handler
	switch contentType {
	case httputils.ContentTypeTextPlain:
		errorHandler = func(params httputils.Params) error {
			return nil
			// e.Text(http.StatusNotFound, "Not Found")
		}
	case httputils.ContentTypeTextHtml:
		errorHandler = func(params httputils.Params) error {
			return nil
		}
	default:
		errorHandler = func(params httputils.Params) error {
			return nil
		}
	}
	return errorHandler
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	handleFuncParams := httputils.Params{ResponseWriter: res, Request: req}
	for _, route := range r.Routes.Data() {
		if matches := route.Pattern.FindStringSubmatch(req.URL.Path); len(matches) > 0 && route.Method == req.Method {
			if len(matches) > 1 {
				handleFuncParams.Params = matches[1:]
			}
			route.Handler(handleFuncParams)
			return
		}
	}
	r.Handler(handleFuncParams)
}

/*
// ハンドラーが持てばよいのでは？
func (e *Exchange) Text(code int, body string) {
	e.ResponseWriter.Header().Set("Content-Type", httputils.ContentTypeTextPlain)
	e.WriteHeader(code)
	io.WriteString(e.ResponseWriter, fmt.Sprintf("%s\n", body))
}
*/