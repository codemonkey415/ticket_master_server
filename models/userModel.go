package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User is the model that governs all notes objects retrived or inserted into the DB
type User struct {
	ID                   primitive.ObjectID   `bson:"_id"`
	First_name           *string              `json:"first_name" validate:"required,min=2,max=100"`
	Last_name            *string              `json:"last_name" validate:"required,min=2,max=100"`
	Password             *string              `json:"Password" validate:"required,min=6"`
	Email                *string              `json:"email" validate:"email,required"`
	Token                *string              `json:"token"`
	Refresh_token        *string              `json:"refresh_token"`
	Reset_Password_Token *string              `json:"reset_password_token"`
	Created_at           time.Time            `json:"created_at"`
	Updated_at           time.Time            `json:"updated_at"`
	User_id              string               `json:"user_id"`
	Due_Date             time.Time            `json:"due_date"`
	Is_Approved          bool                 `json:"is_approved"`
	Role                 string               `json:"role"`
	Reservations         []primitive.ObjectID `json:"reservations"`
	// Gender        *string            `json:"gender" validate:"required"`
	// DateOfBirth   *string            `json:"dateOfBirth" validate:"required"`
}
