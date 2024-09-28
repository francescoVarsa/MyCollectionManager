package V1

import (
	"smart_modellism/pkg/handlers"
	"smart_modellism/router"
)

type RoutesV1 struct {
	Routes router.Routes
}

var Endpoints RoutesV1

func (r *RoutesV1) Get() router.Routes {
	routes := router.Routes{{
		Path:    "/api/v1/models",
		Method:  "GET",
		Handler: handlers.GetModels,
	}, {
		Path:    "/api/v1/models/:id",
		Method:  "GET",
		Handler: handlers.GetModelById,
	}, {
		Path:    "/api/v1/models/create",
		Method:  "POST",
		Handler: handlers.InsertModel,
	},
	}

	return routes
}
