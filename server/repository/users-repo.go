package repository

import "../models"

type UsersRepository interface {
	AddUser(user *models.User) (*models.User, error) // InsertOne
}