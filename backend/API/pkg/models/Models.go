package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ExternalLinks struct {
	Pages []string `json:"pages"`
	Media struct {
		Images []string `json:"images"`
		Videos []string `json:"videos"`
	} `json:"media"`
}

type NewModel struct {
	ModelName     string        `json:"modelName" bson:"modelName" binding:"required"`
	Category      string        `json:"modelType" bson:"modelType" binding:"required"`
	Scale         string        `json:"scale" bson:"scale"`
	Material      string        `json:"material" bson:"material"`
	Manufacturer  string        `json:"manufacturer" bson:"manufacturer"`
	Images        []string      `json:"images" bson:"images"`
	Tags          []string      `json:"tags" bson:"tags"`
	BuidingTime   int           `json:"buildingTime" bson:"buildingTime"`
	Difficulty    int           `json:"difficulty" bson:"difficulty"`
	ExternalLinks ExternalLinks `json:"links" bson:"links"`
	Description   string        `json:"description" bson:"description" binding:"required"`
	CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
}

type Model struct {
	Id            primitive.ObjectID `json:"id" bson:"_id" binding:"required"`
	ModelName     string             `json:"modelName" bson:"modelName" binding:"required"`
	Category      string             `json:"modelType" bson:"modelType" binding:"required"`
	Scale         string             `json:"scale" bson:"scale"`
	Material      string             `json:"material" bson:"material"`
	Manufacturer  string             `json:"manufacturer" bson:"manufacturer"`
	Images        []string           `json:"images" bson:"images"`
	Tags          []string           `json:"tags" bson:"tags"`
	BuidingTime   int                `json:"buildingTime" bson:"buildingTime"`
	Difficulty    int                `json:"difficulty" bson:"difficulty"`
	ExternalLinks ExternalLinks      `json:"links" bson:"links"`
	Description   string             `json:"description" bson:"description" binding:"required"`
	CreatedAt     time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time          `json:"updated_at" bson:"updated_at"`
}

type ModelUpdate struct {
	Id            *string        `json:"id" bson:"_id" binding:"required"`
	ModelName     *string        `json:"modelName,omitempty"`
	Category      *string        `json:"modelType,omitempty"`
	Scale         *string        `json:"scale,omitempty"`
	Material      *string        `json:"material,omitempty"`
	Manufacturer  *string        `json:"manufacturer,omitempty"`
	Images        *[]string      `json:"images,omitempty"`
	Tags          *[]string      `json:"tags,omitempty"`
	BuidingTime   *int           `json:"buildingTime,omitempty"`
	Difficulty    *int           `json:"difficulty,omitempty"`
	ExternalLinks *ExternalLinks `json:"links,omitempty"`
	Description   *string        `json:"description,omitempty"`
	UpdatedAt     time.Time      `json:"-"`
}

type ModelsDocuments []Model
