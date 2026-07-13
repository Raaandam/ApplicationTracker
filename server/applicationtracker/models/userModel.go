package model

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
	ID           bson.ObjectID `bson:"_id" json:"id"`
	Name         string        `bson:"name" json:"name" validate:"required"`
	Email        string        `bson:"email" json:"email" validate:"required"`
	Password     string        `bson:"password" json:"password" validate:"required,min=6"`
	CreatedAt    time.Time     `bson:"created_at" json:"created_at" validate:"required"`
	Token        *string       `bson:"token" json:"token"`
	RefreshToken *string       `bson:"refresh_token" json:"refresh_token"`
	UpdatedAt    time.Time     `bson:"updated_at" json:"updated_at"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UserResponse struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	//not goint to implement token and refresh token. will send http only cookies using c.SetCookie()
}
