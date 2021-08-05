package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	Name string `bson:"name,omitempty"`
	DOB	string `bson:"dob,omitempty"`
	Address string `bson:"address,omitempty"`
	Description string `bson:"description,omitempty"`
	CreatedAt string `bson:"createdAt,omitempty"`
	Followers []string `bson:"followers,omitempty"`
	Following []string `bson:"following,omitempty"`
	Latitude float64 `bson:"latitude,omitempty"`
	Longitude float64 `bson:"longitude,omitempty"`
}
