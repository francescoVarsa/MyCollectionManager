package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"smart_modellism/pkg/database"
	"smart_modellism/pkg/models"
	"smart_modellism/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetModels(ctx *gin.Context) {

	res, err := database.ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			return DB.Database("admin").Collection("cars_models").Find(
				context.TODO(),
				bson.D{},
			)
		},
	)

	if err != nil {
		panic(err)
	}

	if result, ok := res.(*mongo.Cursor); ok {
		var docs models.Models

		if err := result.All(ctx, &docs); err != nil {
			utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

			return
		}

		var responseJSON utils.SuccessResponse
		responseJSON.Status = "Success"
		responseJSON.Data = docs

		ctx.JSON(http.StatusOK, responseJSON)
	}

}

func GetModelById(ctx *gin.Context) {
	id, exists := ctx.Params.Get("id")

	if !exists {
		utils.ErrorJSON(errors.New("invalid request: missing param id"), ctx, http.StatusBadRequest)

		return
	}

	mongoId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		utils.ErrorJSON(errors.New("cannot retrieve the mongo object Id"), ctx, http.StatusInternalServerError)

		return
	}

	res, _ := database.ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			return DB.Database("admin").Collection("cars_models").FindOne(context.Background(), bson.M{
				"_id": mongoId,
			}), nil
		},
	)

	result, ok := res.(*mongo.SingleResult)
	var document models.Model

	if ok {
		if err := result.Decode(&document); err != nil {
			return
		}

		var responseJSON utils.SuccessResponse

		responseJSON.Status = "Success"
		responseJSON.Data = document

		ctx.JSON(http.StatusOK, responseJSON)
	}
}

func InsertModel(ctx *gin.Context) {
	var resource models.NewModel

	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		utils.ErrorJSON(err, ctx, 500)

		return
	}

	err = json.Unmarshal(body, &resource)

	if err != nil {
		utils.ErrorJSON(err, ctx, 500)

		return
	}

	resource.CreatedAt = time.Now()
	resource.UpdatedAt = time.Now()

	res, err := database.ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			return DB.Database("admin").Collection("cars_models").InsertOne(context.Background(), resource)
		},
	)

	if err != nil {
		utils.ErrorJSON(err, ctx, 500)

		return
	}

	if result, ok := res.(*mongo.InsertOneResult); ok {
		newDocId, _ := result.InsertedID.(primitive.ObjectID)

		res, _ := database.ExecuteQuery(
			func(DB *mongo.Client) (interface{}, error) {
				return DB.Database("admin").Collection("cars_models").FindOne(context.Background(), bson.M{"_id": newDocId}), nil
			},
		)

		if res, ok := res.(*mongo.SingleResult); ok {
			var document models.Model

			if err := res.Decode(&document); err != nil {
				return
			}

			var responseJSON utils.SuccessResponse

			responseJSON.Status = "Success"
			responseJSON.Data = document

			ctx.JSON(http.StatusCreated, responseJSON)

		}

	}
}
