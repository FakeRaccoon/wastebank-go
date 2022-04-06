package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username"`
	Name     string             `json:"name"`
	Balance  int                `json:"balance"`
	Access   AccessStruct       `json:"access"`
}

type UserWithPassword struct {
	Id       primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username string             `json:"username"`
	Name     string             `json:"name"`
	Password string             `json:"password"`
	Access   AccessStruct       `json:"access"`
}

type AccessStruct struct {
	Role string `json:"role"`
}
