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
	result := database.ReadAll(ctx, "models", "documents")

	if result == nil {
		return
	}

	var docs models.ModelsDocuments

	if err := result.All(ctx, &docs); err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return
	}

	var responseJSON utils.SuccessResponse
	responseJSON.Status = "Success"
	responseJSON.Data = docs

	ctx.JSON(http.StatusOK, responseJSON)

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

	filter := bson.M{
		"_id": mongoId,
	}

	result := database.ReadOne(ctx, filter, "models", "documents")

	var document models.Model

	if ok := database.DecodeSingleResult(ctx, result, &document); ok {
		var responseJSON utils.SuccessResponse

		responseJSON.Status = "Success"
		responseJSON.Data = document

		ctx.JSON(http.StatusOK, responseJSON)
	}

}

func InsertModel(ctx *gin.Context) {
	var resource models.NewModel
	const dbName, collectionName = "models", "documents"

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

	res := database.Create(ctx, resource, dbName, collectionName)

	if res == nil {
		return
	}

	newDocId := res.InsertedID

	filter := bson.M{"_id": newDocId}
	var document models.Model

	result := database.ReadOne(ctx, filter, dbName, collectionName)

	if ok := database.DecodeSingleResult(ctx, result, &document); ok {
		var responseJSON utils.SuccessResponse

		responseJSON.Status = "Success"
		responseJSON.Data = document

		ctx.JSON(http.StatusCreated, responseJSON)
	}
}

func DeleteModel(ctx *gin.Context) {
	id, exists := ctx.Params.Get("id")

	if !exists {
		utils.ErrorJSON(errors.New("invalid request: missing param id"), ctx, http.StatusBadRequest)

		return
	}

	res, err := database.ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {

			objectID, err := primitive.ObjectIDFromHex(id)

			if err != nil {
				return nil, err
			}

			return DB.Database("admin").Collection("cars_models").DeleteOne(context.Background(), bson.M{"_id": objectID})
		},
	)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return
	}

	if deleteResult, ok := res.(*mongo.DeleteResult); ok {
		itemsDeleted := deleteResult.DeletedCount

		var responseJSON utils.SuccessResponse

		type deleteDataResponse struct {
			Count     int64  `json:"deletedElements"`
			DeletedId string `json:"deletedElementId"`
		}

		dataResponse := deleteDataResponse{
			Count:     itemsDeleted,
			DeletedId: id,
		}

		responseJSON.Status = "Success"
		responseJSON.Data = dataResponse

		ctx.JSON(http.StatusAccepted, responseJSON)

		return
	}

	utils.ErrorJSON(errors.New("the document has been deleted but something goes wrong generating a response"), ctx, http.StatusInternalServerError)
}
