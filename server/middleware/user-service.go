package middleware

import (
	"errors"
	"time"

	"../models"
	"../repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// logrus "github.com/sirupsen/logrus"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]primitive.M, error)
	GetUserById(user string) (bson.M, error)
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

func (*service) GetAllUsers() ([]primitive.M, error) {
	users, err := repo.GetAllUsers()
	if err != nil {
		return nil, err
	} else {
		return users, nil
	}
}

func (*service) GetUserById(user string) (bson.M, error) {
	userResult, err := repo.GetUser(user)
	if err != nil {
		return nil, err
	} else {
		return userResult, nil
	}
}
