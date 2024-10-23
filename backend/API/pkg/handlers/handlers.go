package handlers

import (
	"errors"
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
	id, ok := getIdParam(ctx)

	if !ok {
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

	err := decodeRequestBody(ctx, &resource)

	if err != nil {
		utils.ErrorJSON(err, ctx, 500)

		return
	}

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

	id, ok := getIdParam(ctx)

	if !ok {
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

func ModifyModel(ctx *gin.Context) {
	var resource models.ModelUpdate
	const dbName, collectionName = "models", "documents"

	err := decodeRequestBody(ctx, &resource)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return
	}

	filter := bson.M{}
	updateMap := bson.M{}

	var objectId primitive.ObjectID

	if resource.Id != nil {
		objectId, err = primitive.ObjectIDFromHex(*resource.Id)

		if err != nil {
			utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

			return
		}

		filter["_id"] = objectId
	} else {
		utils.ErrorJSON(errors.New("missing id field in the request"), ctx, http.StatusBadRequest)

		return
	}

	if resource.ModelName != nil {
		updateMap["modelName"] = resource.ModelName
	}

	if resource.Category != nil {
		updateMap["modelType"] = resource.Category
	}

	if resource.Scale != nil {
		updateMap["scale"] = resource.Scale
	}

	if resource.Material != nil {
		updateMap["material"] = resource.Material
	}

	if resource.Manufacturer != nil {
		updateMap["manufacturer"] = resource.Manufacturer
	}

	if resource.Images != nil {
		updateMap["images"] = resource.Images
	}

	if resource.Tags != nil {
		updateMap["tags"] = resource.Tags
	}

	if resource.BuidingTime != nil {
		updateMap["buildingTime"] = resource.BuidingTime
	}

	if resource.Difficulty != nil {
		updateMap["difficulty"] = resource.Difficulty
	}

	if resource.ExternalLinks != nil {
		updateMap["links"] = resource.ExternalLinks
	}

	if resource.Description != nil {
		updateMap["description"] = resource.Description
	}

	resource.UpdatedAt = time.Now()

	res, ok := database.Update(ctx, filter, updateMap, dbName, collectionName)

	sendUpdateResponse(ctx, res, ok, objectId, dbName, collectionName)
}

func ReplaceModel(ctx *gin.Context) {
	var replacement models.Model
	const dbName, collectionName = "models", "documents"

	err := decodeRequestBody(ctx, &replacement)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return
	}

	res, ok := database.Replace(ctx, replacement, dbName, collectionName)

	sendUpdateResponse(ctx, res, ok, replacement.Id, dbName, collectionName)

}
