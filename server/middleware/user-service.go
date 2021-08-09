package middleware

import (
	"errors"
	"net/http"
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
	DeleteUserById(user string) error
	UpdateUserById(userId string, r *http.Request) error
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

func (*service) DeleteUserById(user string) error {
	err := repo.DeleteUser(user)
	return err
}

func (*service) UpdateUserById(userId string, r *http.Request) error {
	name := r.URL.Query().Get("name")
	dob := r.URL.Query().Get("dob")
	address := r.URL.Query().Get("address")
	description := r.URL.Query().Get("description")
	latitude := r.URL.Query().Get("latitude")
	longitude := r.URL.Query().Get("longitude")
	change := bson.M{}
	if name != "" {
		change["name"] = name
	}
	if dob != "" {
		change["dob"] = dob
	}
	if address != "" {
		change["address"] = address
	}
	if description != "" {
		change["description"] = description
	}
	if latitude != "" {
		change["latitude"] = latitude
	}
	if longitude != "" {
		change["longitude"] = longitude
	}
	update := bson.M{"$set": change}
	err := repo.UpdateUser(userId, update)
	return err
}
