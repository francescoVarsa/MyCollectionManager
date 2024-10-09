package database

import (
	"context"
	"errors"
	"net/http"
	"smart_modellism/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbQueryExecutionMaxTime = time.Millisecond * 80

func Create(
	ctx *gin.Context, resource interface{}, dbName string, dbCollection string) *mongo.InsertOneResult {
	res, err := ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			return DB.Database(dbName).Collection(dbCollection).InsertOne(context.Background(), resource)
		},
	)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return nil
	}

	if res, ok := res.(*mongo.InsertOneResult); ok {
		return res
	}

	handleUnexpectedMongoType(ctx)

	return nil
}

func Read(
	ctx *gin.Context, filter interface{}, dbName string, dbCollection string) *mongo.Cursor {

	res, err := ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			context, cancel := context.WithTimeout(context.Background(), time.Duration(dbQueryExecutionMaxTime))
			defer cancel()

			return DB.Database(dbName).Collection(dbCollection).Find(
				context,
				filter,
			)
		},
	)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return nil
	}

	if result, ok := res.(*mongo.Cursor); ok {
		return result
	}

	handleUnexpectedMongoType(ctx)

	return nil
}

func ReadAll(ctx *gin.Context, dbName string, dbCollection string) *mongo.Cursor {
	return Read(ctx, bson.D{}, dbName, dbCollection)
}

func ReadOne(ctx *gin.Context, filter interface{}, dbName string, dbCollection string) *mongo.SingleResult {
	res, _ := ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			return DB.Database(dbName).Collection(dbCollection).FindOne(context.Background(), filter), nil
		},
	)

	result, ok := res.(*mongo.SingleResult)

	if ok {
		return result

	}

	handleUnexpectedMongoType(ctx)

	return nil
}

func Update(
	ctx *gin.Context, filter interface{}, dbName string, dbCollection string) {
}

func Delete(
	ctx *gin.Context, filter interface{}, dbName string, dbCollection string) {
}

func handleUnexpectedMongoType(ctx *gin.Context) {
	utils.ErrorJSON(errors.New("unexpected return type"), ctx, http.StatusInternalServerError)
}

func DecodeSingleResult(ctx *gin.Context, result *mongo.SingleResult, v interface{}) bool {
	err := result.Decode(v)

	if err == mongo.ErrNoDocuments {
		utils.ErrorJSON(errors.New("the resource you're looking for does not exists in the database"), ctx, http.StatusNotFound)

		return false
	} else if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return false
	}

	return true
}
