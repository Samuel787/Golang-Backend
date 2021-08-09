package repository

import (
	"../models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UsersRepository interface {
	AddUser(user *models.User) (*models.User, error)
	GetAllUsers() ([]primitive.M, error)
	GetUser(userId string) (bson.M, error)
}
