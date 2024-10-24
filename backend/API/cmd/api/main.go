package main

import (
	"flag"
	"log"
	"os"
	"smart_modellism/pkg/config"
	"smart_modellism/router"
	"smart_modellism/router/routes/V1"

	"github.com/spf13/viper"
)

func main() {
	var r router.Router

	envFilename := flag.String("envFile", "", "path to the .env file")

	// Parso i flag
	flag.Parse()

	// Controlla se il file di configurazione esiste
	if *envFilename == "" {
		log.Fatal("envFile argument not provided")
	} else if _, err := os.Stat(*envFilename); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s\n", *envFilename)

	}

	viper.SetConfigFile(*envFilename)

	r.Init()
	r.SetRoutes(V1.Endpoints.Get())

	port, err := config.GetEnv("SERVER_PORT")

	if err != nil {
		panic(err)
	}

	r.Start(port)
}
