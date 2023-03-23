package pipeline

import (
	"github.com/gorilla/mux"
	"net/http"
)

// Route 路由
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	Middlewares []Handler
}

// AddRoutes 添加路由
// AddRoutes 原生http.ServerMux 功能比较单一，所以这里使用github.com/gorilla/mux
func AddRoutes(routes []Route, router *mux.Router) {
	var (
		handler  http.Handler
		pipeline *Pipeline
	)

	for _, route := range routes {
		if len(route.Middlewares) > 0 {
			pipeline = NewPipeline()

			for _, middleware := range route.Middlewares {
				pipeline.Use(middleware)
			}

			pipeline.Use(Adapt(route.HandlerFunc))
			handler = pipeline
		} else {
			handler = route.HandlerFunc
		}

		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(handler)
	}
}
