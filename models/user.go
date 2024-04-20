package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"  json:"id,omitempty"`
	Email     string             `bson:"email"           json:"email"`
	Password  string             `bson:"password,omitempty" json:"-"`
	FirstName string             `bson:"first_name"       json:"first_name"`
	LastName  string             `bson:"last_name"        json:"last_name"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
