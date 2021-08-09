package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"../models"
	"../repository"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	jwt "github.com/dgrijalva/jwt-go"
	logrus "github.com/sirupsen/logrus"
)

var mySigningKey = []byte("mysupersecretphrase")

const connectionString = "mongodb+srv://SamuelRyde:SlfZ0ehN2bDKnr3h@rydecluster.pbok3.mongodb.net/test?retryWrites=true&w=majority"
const dbName = "rydedb"
const collectionName = "users"

var collection *mongo.Collection
var ApiService UserService

func init() {
	enableLogging(true)
	ApiService = NewUserService(repository.NewMongoUsersRepository())
}

func enableLogging(flag bool) {
	if flag {
		logrus.SetOutput(os.Stdout)
		logrus.SetFormatter(&logrus.JSONFormatter{})
		logLevel, err := logrus.ParseLevel("debug")
		if err != nil {
			logLevel = logrus.InfoLevel
		}
		logrus.SetLevel(logLevel)
	}
}

func AuthorizeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There is an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"auth-msg": "Error during token parsing",
					"token":    token,
					"error":    err.Error(),
					"time":     time.Now().String(),
				}).Warn("AuthorizeUser")
			}
			if token.Valid {
				logrus.WithFields(logrus.Fields{
					"auth-msg": "Good token. Allowing access to API",
					"token":    token,
					"time":     time.Now().String(),
				}).Info("AuthorizeUser")
				next.ServeHTTP(w, r)
			}
		} else {
			logrus.WithFields(logrus.Fields{
				"auth-msg": "An attempt to use API without auth token",
				"time":     time.Now().String(),
			}).Warn("AuthorizeUser")
			json.NewEncoder(w).Encode("Not Authorized")
		}
	})
}

// Create User route
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	err := ApiService.CreateUser(&user)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user-object": user,
			"error":       err.Error(),
			"time":        time.Now().String(),
		}).Warn("CrateUser")
		json.NewEncoder(w).Encode(err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"user-object": user,
			"msg":         "Successfully created the user",
			"time":        time.Now().String(),
		}).Info("CrateUser")
		json.NewEncoder(w).Encode("Success")
	}
}

// Get All Users route
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload, err := ApiService.GetAllUsers()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"users-object": payload,
			"error":        err.Error(),
			"time":         time.Now().String(),
		}).Warn("GetAllUsers")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"users": payload,
			"msg":   "Successfully queried for all users in database",
			"time":  time.Now().String(),
		}).Info("GetAllUsers")
		json.NewEncoder(w).Encode(payload)
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	user, err := ApiService.GetUserById(params["id"])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user":  user,
			"error": err.Error(),
			"time":  time.Now().String(),
		}).Warn("GetUser")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"users": user,
			"msg":   "Successfully queried for user in database",
			"time":  time.Now().String(),
		}).Info("GetUser")
		json.NewEncoder(w).Encode(user)
	}
}

func findUserById(user string) bson.M {
	var result bson.M
	id, _ := primitive.ObjectIDFromHex(user)
	collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	return result
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	err := ApiService.DeleteUserById(params["id"])
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user-id": params["id"],
			"error":   err.Error(),
			"time":    time.Now().String(),
		}).Warn("DeleteUser")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"user-id": params["id"],
			"msg":     "Successfully deleted user from database",
			"time":    time.Now().String(),
		}).Info("DeleteUser")
		json.NewEncoder(w).Encode("Success")
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	hex_id := r.URL.Query().Get("id")
	err := ApiService.UpdateUserById(hex_id, r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"user-id": hex_id,
			"error":   err.Error(),
			"time":    time.Now().String(),
		}).Warn("UpdateUser")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"user-id": hex_id,
			"msg":     "Successfully updated user in database",
			"time":    time.Now().String(),
		}).Info("UpdateUser")
		json.NewEncoder(w).Encode("Success")
	}
}

/**
Adds followers to the user
Ensures that user is not adding himself as follower
Ensures that user is not attempting to add non-existent user as follower
*/
func AddFollowerToUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userId := r.URL.Query().Get("userId")
	followerId := r.URL.Query().Get("followerId")
	err := ApiService.AddFollowerToUser(userId, followerId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userId":     userId,
			"followerId": followerId,
			"error":      err.Error(),
			"time":       time.Now().String(),
		}).Warn("AddFollowerToUser")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"userId":     userId,
			"followerId": followerId,
			"msg":        "Successfully added follower to user",
			"time":       time.Now().String(),
		}).Info("AddFollowerToUser")
		json.NewEncoder(w).Encode("Success")
	}
}

/**
Adds following to the user
Ensures that user is not adding himself as following
Ensures that user is not attempting to add non-existent user as following
@TODO: Test this method
*/
func AddFollowingToUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userId := r.URL.Query().Get("userId")
	followingId := r.URL.Query().Get("followingId")
	err := ApiService.AddFollowingToUser(userId, followingId)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userId":      userId,
			"followingId": followingId,
			"error":       err.Error(),
			"time":        time.Now().String(),
		}).Warn("AddFollowingToUser")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"userId":      userId,
			"followingId": followingId,
			"msg":         "Successfully added following to user",
			"time":        time.Now().String(),
		}).Info("AddFollowingToUser")
		json.NewEncoder(w).Encode("Success")
	}
}

/**
@TODO: Test this method
User will be able to get the list of users who they are *following* who are nearby.
This will not include users who are following that user (unless this user is following them as well)
because of a directed relationship. user knows about users he is following but not about
users who are following him but did not let this user follow them back.
^ Explain and document the above more clearly
*/
func GetNearByFollowing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	userId := r.URL.Query().Get("userId")
	dist := r.URL.Query().Get("dist")
	limit := r.URL.Query().Get("limit")
	results, err := ApiService.GetNearByUsers(userId, dist, limit)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"userId": userId,
			"error":  err.Error(),
			"time":   time.Now().String(),
		}).Warn("GetNearByFollowing")
		json.NewEncoder(w).Encode("Error: " + err.Error())
	} else {
		logrus.WithFields(logrus.Fields{
			"results": results,
			"msg":     "Successfully retrieved nearby users",
			"time":    time.Now().String(),
		}).Info("GetNearByFollowing")
		json.NewEncoder(w).Encode(results)
	}
}
