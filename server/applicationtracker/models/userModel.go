package model

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID        bson.ObjectID  `bson:"_id" json:"_id" validate:"required"`
	Name      string         `bson:"name" json:"name" validate:"required"`
	Email     string         `bson:"email" json:"email" validate:"required"`
	Password  string         `bson:"password" json:"password" validate:"required"`
	CreatedAt bson.Timestamp `bson:"created_at" json:"created_at" validate:"required"`
}
