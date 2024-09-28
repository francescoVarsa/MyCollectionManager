package handlers

import (
	"context"
	"net/http"
	"smart_modellism/pkg/database"
	"smart_modellism/pkg/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetModels(ctx *gin.Context) {

	res, err := database.ExecuteQuery(
		func(DB *mongo.Client) (*mongo.Cursor, error) {
			return DB.Database("admin").Collection("cars_models").Find(
				context.TODO(),
				bson.D{},
			)
		},
	)

	if err != nil {
		panic(err)
	}

	var docs []models.Models

	if err := res.All(ctx, &docs); err != nil {
		panic(err)
	}

	response := struct {
		Data []models.Models
	}{
		Data: docs,
	}

	ctx.JSON(http.StatusOK, response)
}

func GetModelById(ctx *gin.Context) {}

func InsertModel(ctx *gin.Context) {}
