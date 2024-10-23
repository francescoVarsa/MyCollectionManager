package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"smart_modellism/pkg/database"
	"smart_modellism/pkg/models"
	"smart_modellism/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getIdParam(ctx *gin.Context) (id string, ok bool) {
	id, exists := ctx.Params.Get("id")

	if !exists {
		utils.ErrorJSON(errors.New("invalid request: missing param id"), ctx, http.StatusBadRequest)

		return "", false
	}

	return id, true
}

func decodeRequestBody(ctx *gin.Context, resource interface{}) error {
	body, err := io.ReadAll(ctx.Request.Body)

	if err != nil {
		return err
	}

	err = json.Unmarshal(body, resource)

	return err
}

func sendUpdateResponse(ctx *gin.Context, res *mongo.UpdateResult, ok bool, objId primitive.ObjectID, dbName string, collectionName string) {
	if !ok {
		return
	}

	if res.MatchedCount == 0 {
		utils.ErrorJSON(errors.New("the document you are trying to modify doesn't exists"), ctx, http.StatusNotFound)

		return
	}

	filter := bson.M{"_id": objId}

	result, ok := database.ReadOne(ctx, filter, dbName, collectionName)

	if !ok {
		return
	}

	var document models.Model

	if ok := database.DecodeSingleResult(ctx, result, &document); ok {
		var responseJSON utils.SuccessResponse

		responseJSON.Status = "Success"
		responseJSON.Data = document

		codeStatus := http.StatusOK

		if res.ModifiedCount == 0 {
			codeStatus = http.StatusNotModified
		}
		ctx.JSON(codeStatus, responseJSON)

		return
	}

	utils.ErrorJSON(errors.New("the resource was not modified"), ctx, http.StatusInternalServerError)
}
