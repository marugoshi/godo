package router

import (
	"context"
	"github.com/marugoshi/gobm/presentation/handler"
	"github.com/marugoshi/gobm/presentation/httputils"
	"net/http"
	"regexp"
)

func (r *Router) routesData() []Route {
	return []Route{
		Route{`^/bookmarks$`, http.MethodGet, r.Handlers.Bookmarks},
		Route{`^/bookmarks/(\d*)$`, http.MethodGet, r.Handlers.Bookmark},
	}
}

type Route struct {
	Pattern string
	Method  string
	httputils.Func
}

type Router struct {
	handler.Handlers
	httputils.Func
}

func NewRouter(contentType string) *Router {
	return &Router{handler.NewHandlers(), notFoundError(contentType)}
}

func notFoundError(contentType string) httputils.Func {
	var errorHandler httputils.Func
	switch contentType {
	case httputils.ContentTypeTextPlain:
		errorHandler = notFoundErrorTextPlain
	case httputils.ContentTypeTextHtml:
		errorHandler = notFoundErrorTextHtml
	default:
		errorHandler = notFoundErrorTextPlain
	}
	return errorHandler
}

func notFoundErrorTextPlain(ctx context.Context, api httputils.Api) error {
	return nil
}

func notFoundErrorTextHtml(ctx context.Context, api httputils.Api) error {
	return nil
}

func (r *Router) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	api := httputils.Api{ResponseWriter: res, Request: req}
	for _, route := range r.routesData() {
		re := regexp.MustCompile(route.Pattern)
		if matches := re.FindStringSubmatch(req.URL.Path); len(matches) > 0 && route.Method == req.Method {
			if len(matches) > 1 {
				api.Params = matches[1:]
			}
			err := route.Func(ctx, api)
			if err != nil {
				return
			} else {
				// TODO: error
				return
			}
		}
	}
	r.Func(ctx, api)
}
