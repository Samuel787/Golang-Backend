package middleware

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"../models"
	"../repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserService interface {
	CreateUser(user *models.User) error
	GetAllUsers() ([]primitive.M, error)
	GetUserById(user string) (bson.M, error)
	DeleteUserById(user string) error
	UpdateUserById(userId string, r *http.Request) error
	AddFollowerToUser(userId string, followerId string) error
	AddFollowingToUser(userId string, followingId string) error
	GetNearByUsers(userId string, dist string, limit string) ([]primitive.M, error)
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

func (*service) AddFollowerToUser(userId string, followerId string) error {
	if userId == followerId {
		return errors.New("[AddFollowerToUser] User can't add himself as follower")
	}
	user, err := repo.GetUser(userId)
	if err != nil {
		return errors.New("[AddFollowerToUser] the userId does not exist in the database")
	}
	_, err2 := repo.GetUser(followerId)
	if err2 != nil {
		return errors.New("[AddFollowerToUser] the followingId does not exist in the database")
	}
	if user["followers"] == nil {
		user["followers"] = [1]string{followerId}
	} else {
		var existingFollowers bson.A = user["followers"].(bson.A)
		for _, currFollower := range existingFollowers {
			if currFollower == followerId {
				return errors.New("[AddFollowerToUser] User has already added this follower")
			}
		}
		user["followers"] = append(existingFollowers, followerId)
	}
	update := bson.M{"$set": user}
	errUpdate := repo.UpdateUser(userId, update)
	return errUpdate
}

func (*service) AddFollowingToUser(userId string, followingId string) error {
	if userId == followingId {
		return errors.New("[AddFollowerToUser] User can't follow himself")
	}
	user, err := repo.GetUser(userId)
	if err != nil {
		return errors.New("[AddFollowingToUser] the userId does not exist in the database")
	}
	_, err2 := repo.GetUser(followingId)
	if err2 != nil {
		return errors.New("[AddFollowingToUser] the followingId does not exist in the database")
	}
	if user["following"] == nil {
		user["following"] = [1]string{followingId}
	} else {
		var existingFollowing bson.A = user["following"].(bson.A)
		for _, currFollower := range existingFollowing {
			if currFollower == followingId {
				return errors.New("[AddFollowingToUser] The user is already following that user")
			}
		}
		user["following"] = append(existingFollowing, followingId)
	}
	update := bson.M{"$set": user}
	errUpdate := repo.UpdateUser(userId, update)
	return errUpdate
}

func (*service) GetNearByUsers(userId string, dist string, limit string) ([]primitive.M, error) {
	// check if user exists
	user, errUser := repo.GetUser(userId)
	if errUser != nil {
		return nil, errors.New("[GetNearByUsers] the user with userId does not exist")
	}
	// convert distance to float64
	distance, errDistParsing := strconv.ParseFloat(dist, 64)
	if errDistParsing != nil {
		return nil, errors.New("[GetNearByUsers] dist value cannot be parsed to float")
	}
	// convert limit to integer
	limitVal, errLimitParsing := strconv.Atoi(limit)
	if errLimitParsing != nil {
		return nil, errors.New("[GetNearByUsers] limit value cannot be parsed to integer")
	}
	// run the algorithm
	if user["latitude"] == nil || user["longitude"] == nil {
		return nil, errors.New("[GetNearByUsers] user does not have both latitude and longitude values")
	}

	userLat := user["latitude"].(float64)
	userLong := user["longitude"].(float64)
	var results []primitive.M
	var followingList bson.A
	var count = 0
	if user["following"] == nil {
		return results, nil
	} else {
		followingList = user["following"].(bson.A)
		for _, currFollowing := range followingList {
			if count >= limitVal {
				break
			}
			currUser, _ := repo.GetUser(currFollowing.(string))
			if currUser != nil && currUser["latitude"] != nil && currUser["longitude"] != nil {
				var lat = currUser["latitude"].(float64)
				var long = currUser["longitude"].(float64)
				curr_dist := getDist(userLat, userLong, lat, long)
				if curr_dist <= distance {
					results = append(results, currUser)
					count++
				}
			}
		}
	}
	return results, nil
}

/**
helper method to get dist in metres between two location points
https://www.nhc.noaa.gov/gccalc.shtml
*/
func getDist(lat1 float64, long1 float64, lat2 float64, long2 float64) float64 {
	var distX = (lat1 - lat2) * 111000
	var distY = (long1 - long2) * 111000
	var hyptotenuse = math.Sqrt((distX*distX + distY*distY))
	var distMetres = hyptotenuse * 111000
	return distMetres
}
