package router

import "github.com/gin-gonic/gin"

type Router struct {
	Routes Routes
	engine *gin.Engine
}

type Route struct {
	Path    string
	Handler gin.HandlerFunc
	Method  string
}

type Routes []Route

func (r *Router) Init() {
	router := gin.Default()
	r.engine = router
}

func (r *Router) Start(port string) {
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
			r.engine.GET(route.Path, route.Handler)
		case "POST":
			r.engine.POST(route.Path, route.Handler)
		case "PUT":
			r.engine.POST(route.Path, route.Handler)
		case "PATCH":
			r.engine.POST(route.Path, route.Handler)
		case "DELETE":
			r.engine.POST(route.Path, route.Handler)
		}
	}

}
