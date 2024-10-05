package database

import (
	"context"
	"fmt"
	"smart_modellism/pkg/config"
	"smart_modellism/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type QueryFunc func(DB *mongo.Client) (interface{}, error)

func Connect(uri string) *mongo.Client {
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a new client and connect to the server
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		panic(err)
	}

	// Send a ping to confirm a successful connection
	var result bson.M
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Decode(&result); err != nil {
		panic(err)
	}

	return client
}

func Close(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Disconnect(ctx); err != nil {
		panic(err)
	}
}

func ExecuteQuery(query QueryFunc) (interface{}, error) {
	dbConnString, err := config.GetEnv("DATABASE_CONNECTION_STRING")

	if err != nil {
		panic(err)
	}

	DB := Connect(dbConnString)
	defer Close(DB)

	return query(DB)

}

func HandleFindOneError(err error, id string, ctx *gin.Context) error {
	if err != nil {
		if err == mongo.ErrNoDocuments {
			customError := fmt.Errorf("no document exists with this id: %s", id)

			utils.ErrorJSON(customError, ctx, 404)

		} else {
			utils.ErrorJSON(err, ctx, 500)
		}

		return err
	}

	return nil
}
