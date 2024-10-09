package router

import (
	"smart_modellism/pkg/throttler"
	"smart_modellism/router/routes/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

type Router struct {
	Routes Routes
	engine *gin.Engine
}

type Route struct {
	Path       string
	Handler    gin.HandlerFunc
	Method     string
	Middleware []gin.HandlerFunc
}

type Routes []Route

func (r *Router) Init() {
	router := gin.Default()
	r.engine = router
}

func (r *Router) Start(port string) {
	var throttler throttler.ThrottleRequests
	throttler.SetDuration(1 * time.Second)
	throttler.SetLimit(10)

	go throttler.ClientsCleanup(3 * time.Second)

	r.engine.Use(middleware.RateLimiterMiddleware(throttler))
	r.registerRoutes()
	r.engine.Run(port)
}

func (r *Router) SetRoutes(routes Routes) {
	r.Routes = routes
}

func (r *Router) registerRoutes() {
	for _, route := range r.Routes {
		switch route.Method {
		case "GET":
			r.engine.GET(route.Path, implementHandler(route)...)
		case "POST":
			r.engine.POST(route.Path, implementHandler(route)...)
		case "PUT":
			r.engine.PUT(route.Path, implementHandler(route)...)
		case "PATCH":
			r.engine.PATCH(route.Path, implementHandler(route)...)
		case "DELETE":
			r.engine.DELETE(route.Path, implementHandler(route)...)
		}
	}
}

func implementHandler(r Route) []gin.HandlerFunc {
	var hasMiddleware bool

	if len(r.Middleware) > 0 {
		for _, m := range r.Middleware {
			hasMiddleware = false

			if m != nil {
				hasMiddleware = true

				break
			}
		}
	}

	var handler []gin.HandlerFunc

	if !hasMiddleware {
		return append(handler, r.Handler)
	}

	return append(r.Middleware, r.Handler)
}
