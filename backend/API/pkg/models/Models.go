package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id            primitive.ObjectID `bson:"_id"`
	ModelName     string             `json:"modelName"`
	Category      string             `json:"modelType"`
	Scale         string             `json:"scale"`
	Material      string             `json:"material"`
	Manufacturer  string             `json:"manufacturer"`
	Images        []string           `json:"images"`
	Tags          []string           `json:"tags"`
	BuidingTime   int                `json:"buildingTime"`
	Difficulty    int                `json:"difficulty"`
	ExternalLinks struct {
		Pages []string `json:"pages"`
		Media struct {
			Images []string `json:"images"`
			Videos []string `json:"videos"`
		} `json:"media"`
	} `json:"links"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Models []Model
