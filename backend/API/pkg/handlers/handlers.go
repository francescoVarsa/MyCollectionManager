package handlers

import (
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
)

func GetModels(ctx *gin.Context) {
	result, ok := database.ReadAll(ctx, "models", "documents")

	if !ok {
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

	result, ok := database.ReadOne(ctx, filter, "models", "documents")

	if !ok {
		return
	}

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

	res, ok := database.Create(ctx, resource, dbName, collectionName)

	if !ok {
		return
	}

	newDocId := res.InsertedID

	filter := bson.M{"_id": newDocId}
	var document models.Model

	result, ok := database.ReadOne(ctx, filter, dbName, collectionName)

	if !ok {
		return
	}

	if ok := database.DecodeSingleResult(ctx, result, &document); ok {
		var responseJSON utils.SuccessResponse

		responseJSON.Status = "Success"
		responseJSON.Data = document

		ctx.JSON(http.StatusCreated, responseJSON)
	}
}

func DeleteModel(ctx *gin.Context) {
	const dbName, collectionName = "models", "documents"
	id, exists := ctx.Params.Get("id")

	if !exists {
		utils.ErrorJSON(errors.New("invalid request: missing param id"), ctx, http.StatusBadRequest)

		return
	}

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return
	}

	filter := bson.M{"_id": objectID}

	deleteResult, ok := database.Delete(ctx, filter, dbName, collectionName)

	if !ok {
		return
	}

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
}
