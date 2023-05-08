package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	ImageUrl    string             `bson:"image" json:"image"`
	Price       float64            `bson:"price" json:"price"`
	Rate        float64            `bson:"rate" json:"rate"`
	Location    string             `bson:"location" json:"location"`
	Stock       int                `bson:"stock" json:"stock"`
	Store       string             `bson:"store" json:"store"`
	Created_at  time.Time          `bson:"created_at" json:"created_at"`
	Updated_at  time.Time          `bson:"updated_at" json:"updated_at"`
	Product_id  string             `bson:"product_id" json:"product_id"`
}
