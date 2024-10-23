package database

import (
	"context"
	"errors"
	"net/http"
	"smart_modellism/pkg/models"
	"smart_modellism/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const dbQueryExecutionMaxTime = time.Millisecond * 80

func Create(
	ctx *gin.Context, resource interface{}, dbName string, dbCollection string) (*mongo.InsertOneResult, bool) {
	res, err := ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			return DB.Database(dbName).Collection(dbCollection).InsertOne(context.Background(), resource)
		},
	)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return nil, false
	}

	if res, ok := res.(*mongo.InsertOneResult); ok {
		return res, true
	}

	handleUnexpectedMongoType(ctx)

	return nil, false
}

func Read(
	ctx *gin.Context, filter interface{}, dbName string, dbCollection string) (*mongo.Cursor, bool) {

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

		return nil, false
	}

	if result, ok := res.(*mongo.Cursor); ok {
		return result, true
	}

	handleUnexpectedMongoType(ctx)

	return nil, false
}

func ReadAll(ctx *gin.Context, dbName string, dbCollection string) (*mongo.Cursor, bool) {
	return Read(ctx, bson.D{}, dbName, dbCollection)
}

func ReadOne(ctx *gin.Context, filter interface{}, dbName string, dbCollection string) (*mongo.SingleResult, bool) {
	res, _ := ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			context, cancel := context.WithTimeout(context.Background(), time.Duration(dbQueryExecutionMaxTime))
			defer cancel()

			return DB.Database(dbName).Collection(dbCollection).FindOne(context, filter), nil
		},
	)

	result, ok := res.(*mongo.SingleResult)

	if ok {
		return result, true

	}

	handleUnexpectedMongoType(ctx)

	return nil, false
}

func Update(
	ctx *gin.Context, filter interface{}, update interface{}, dbName string, dbCollection string) (*mongo.UpdateResult, bool) {

	res, err := ExecuteQuery(func(DB *mongo.Client) (interface{}, error) {
		context, cancel := context.WithTimeout(context.Background(), time.Duration(dbQueryExecutionMaxTime))
		defer cancel()

		updateMap := bson.M{"$set": update}

		return DB.Database(dbName).Collection(dbCollection).UpdateOne(context, filter, updateMap)
	})

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)
	}

	if res, ok := res.(*mongo.UpdateResult); ok {
		return res, true
	}

	handleUnexpectedMongoType(ctx)

	return nil, false
}

func Replace(ctx *gin.Context, resource models.Model, dbName string, dbCollection string) (*mongo.UpdateResult, bool) {

	res, err := ExecuteQuery(func(DB *mongo.Client) (interface{}, error) {
		context, cancel := context.WithTimeout(context.Background(), time.Duration(dbQueryExecutionMaxTime))
		defer cancel()

		id := resource.Id

		filter := bson.M{"_id": id}

		return DB.Database(dbName).Collection(dbCollection).ReplaceOne(context, filter, resource)
	})

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)
	}

	if res, ok := res.(*mongo.UpdateResult); ok {
		return res, true
	}

	handleUnexpectedMongoType(ctx)

	return nil, false
}

func Delete(
	ctx *gin.Context, filter interface{}, dbName string, dbCollection string) (*mongo.DeleteResult, bool) {
	res, err := ExecuteQuery(
		func(DB *mongo.Client) (interface{}, error) {
			context, cancel := context.WithTimeout(context.Background(), time.Duration(dbQueryExecutionMaxTime))
			defer cancel()

			return DB.Database(dbName).Collection(dbCollection).DeleteOne(context, filter)
		},
	)

	if err != nil {
		utils.ErrorJSON(err, ctx, http.StatusInternalServerError)

		return nil, false
	}

	if result, ok := res.(*mongo.DeleteResult); ok {
		return result, true
	}

	handleUnexpectedMongoType(ctx)

	return nil, false
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
