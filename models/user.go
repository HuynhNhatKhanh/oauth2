// models/user.go
package models

import (
	"time"
	"user_login/config"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserCollection is the collection name
var UserCollection *mongo.Collection = config.GetCollection(config.DB, "users")

// User model
type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	OTPLogin        string             `bson:"otp_login,omitempty"`
	Username        string             `json:"username" bson:"username" validate:"required"`
	Password        string             `json:"password" bson:"password" validate:"required"`
	Email           string             `json:"email" bson:"email" validate:"required,email"`
	IsVerified      bool               `json:"is_verified" bson:"is_verified"`
	IsVerifiedLogin bool               `json:"is_verified_login" bson:"is_verified_login"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}
