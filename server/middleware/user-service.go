package middleware

import (
	"errors"
	"time"

	"../models"
	"../repository"
	// logrus "github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user *models.User) error
}

type service struct{}

var (
	repo repository.UsersRepository
)

func NewUserService(repository repository.UsersRepository) UserService {
	repo = repository
	return &service{}
}

func (*service) CreateUser(user *models.User) error {
	if user == nil {
		err := errors.New("[Create User] Unable to create nil user")
		return err
	}
	if user.Name == "" {
		err := errors.New("[Create User] Unable to create user with nil Name")
		return err
	}
	if user.DOB == "" {
		err := errors.New("[Create User] Unable to create user with nil DOB")
		return err
	}
	if user.Address == "" {
		err := errors.New("[Create User] Unable to create user with nil DOB")
		return err
	}
	if user.CreatedAt != "" {
		err := errors.New("[Create User] The creation timing of the user has to be stamped by API")
		return err
	} else {
		user.CreatedAt = time.Now().String()
	}
	if user.Followers != nil {
		err := errors.New("[Create User] Cannot create user with followers")
		return err
	}
	if user.Following != nil {
		err := errors.New("[Create User] Cannot create user with following")
		return err
	}
	if user.Latitude == 0 {
		user.Latitude = 0
	}
	if user.Longitude == 0 {
		user.Longitude = 0
	}
	_, err := repo.AddUser(user)
	return err
}
