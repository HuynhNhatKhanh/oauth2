// models/user.go
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	OTP             string             `bson:"otp,omitempty"`
	OTPLogin        string             `bson:"otp_login,omitempty"`
	Username        string             `json:"username" bson:"username" validate:"required"`
	Password        string             `json:"password" bson:"password" validate:"required"`
	Email           string             `json:"email" bson:"email" validate:"required,email"`
	IsVerified      bool               `json:"is_verified" bson:"is_verified"`
	IsVerifiedLogin bool               `json:"is_verified_login" bson:"is_verified_login"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
}
