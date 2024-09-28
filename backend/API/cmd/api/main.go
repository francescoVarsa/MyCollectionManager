package main

import (
	"smart_modellism/pkg/config"
	"smart_modellism/router"
	"smart_modellism/router/routes/V1"
)

func main() {
	var r router.Router

	r.Init()
	r.SetRoutes(V1.Endpoints.Get())

	port, err := config.GetEnv("SERVER_PORT")

	if err != nil {
		panic(err)
	}

	r.Start(port)
}
