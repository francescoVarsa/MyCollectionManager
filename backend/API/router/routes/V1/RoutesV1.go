package V1

import (
	"smart_modellism/pkg/handlers"
	"smart_modellism/router"
	"smart_modellism/router/routes/middleware"

	"github.com/gin-gonic/gin"
)

type RoutesV1 struct {
	Routes router.Routes
}

var Endpoints RoutesV1

func (r *RoutesV1) Get() router.Routes {
	routes := router.Routes{{
		Path:       "/api/v1/models",
		Method:     "GET",
		Handler:    handlers.GetModels,
		Middleware: []gin.HandlerFunc{},
	}, {
		Path:       "/api/v1/models/:id",
		Method:     "GET",
		Handler:    handlers.GetModelById,
		Middleware: []gin.HandlerFunc{},
	}, {
		Path:    "/api/v1/models/create",
		Method:  "POST",
		Handler: handlers.InsertModel,
		Middleware: []gin.HandlerFunc{
			middleware.ValidateCreateModel(),
		},
	}, {
		Path:       "/api/v1/models/delete/:id",
		Method:     "DELETE",
		Handler:    handlers.DeleteModel,
		Middleware: []gin.HandlerFunc{},
	},
	}

	return routes
}
