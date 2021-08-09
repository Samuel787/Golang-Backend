package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"os"
	"strconv"
	"time"

	"../models"
	"../repository"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	// "go.mongodb.org/mongo-driver/mongo/options"
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
	// clientOptions := options.Client().ApplyURI(connectionString)
	// client, err := mongo.Connect(context.TODO(), clientOptions) // mongo.NewClient(clientOptions)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // check the connection
	// err = client.Ping(context.TODO(), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// database := client.Database(dbName)
	// collection = database.Collection(collectionName)
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
		fmt.Println("This is the auth middleware")

		logrus.WithFields(logrus.Fields{
			"Test": "hello there brown cow",
		}).Info("Auth details")

		if r.Header["Token"] != nil {
			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("There is an error")
				}
				return mySigningKey, nil
			})
			if err != nil {
				fmt.Println(w, err.Error())
			}
			if token.Valid {
				next.ServeHTTP(w, r)
			}
		} else {
			fmt.Println(w, "Not Authorized")
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

// func insertOneUser(user models.User) {
// 	insertResult, err := collection.InsertOne(context.Background(), user)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
// }

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

// Get User by their ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	// payload := findUserById(params["id"])
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

// get all users from the DB and return it
// func getAllUsers() []primitive.M {
// 	cur, err := collection.Find(context.Background(), bson.D{{}})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var results []primitive.M
// 	for cur.Next(context.Background()) {
// 		var result bson.M
// 		e := cur.Decode(&result)
// 		if e != nil {
// 			log.Fatal(e)
// 		}
// 		results = append(results, result)
// 	}
// 	if err := cur.Err(); err != nil {
// 		log.Fatal(err)
// 	}
// 	cur.Close(context.Background())
// 	return results
// }

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

// func deleteOneUser(user string) {
// 	fmt.Println(user)
// 	id, _ := primitive.ObjectIDFromHex(user)
// 	filter := bson.M{"_id": id}
// 	d, err := collection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Deleted Document", d.DeletedCount)
// }

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

// func updateOneUser(r *http.Request) {
// 	hex_id := r.URL.Query().Get("id")
// 	id, _ := primitive.ObjectIDFromHex(hex_id)
// 	filter := bson.M{"_id": bson.M{"$eq": id}}
// 	name := r.URL.Query().Get("name")
// 	dob := r.URL.Query().Get("dob")
// 	address := r.URL.Query().Get("address")
// 	description := r.URL.Query().Get("description")
// 	latitude := r.URL.Query().Get("latitude")
// 	longitude := r.URL.Query().Get("longitude")
// 	change := bson.M{}
// 	if name != "" {
// 		change["name"] = name
// 	}
// 	if dob != "" {
// 		change["dob"] = dob
// 	}
// 	if address != "" {
// 		change["address"] = address
// 	}
// 	if description != "" {
// 		change["description"] = description
// 	}
// 	if latitude != "" {
// 		change["latitude"] = latitude
// 	}
// 	if longitude != "" {
// 		change["longitude"] = longitude
// 	}
// 	update := bson.M{"$set": change}
// 	_, err := collection.UpdateOne(context.Background(), filter, update)
// 	if err != nil {
// 		fmt.Println("UpdateOneUser() result ERROR: ", err)
// 	} else {
// 		fmt.Println("Updated record successfully")
// 	}
// }

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
	// if userId == followerId {
	// 	fmt.Println("User can't add himself as follower")
	// }
	// user := findUserById(userId)
	// follower := findUserById(followerId)
	// if user == nil || follower == nil {
	// 	fmt.Println("One of the user doesn't exist")
	// 	return
	// }
	// var followerExists = false
	// if user["followers"] == nil {
	// 	user["followers"] = [1]string{followerId}
	// } else {
	// 	var existingFollowers bson.A = user["followers"].(bson.A)
	// 	for _, currFollower := range existingFollowers {
	// 		if currFollower == followerId {
	// 			followerExists = true
	// 		}
	// 	}
	// 	if !followerExists {
	// 		user["followers"] = append(existingFollowers, followerId)
	// 		fmt.Println("Follower doesn't exist, hence proceeding to add")
	// 	} else {
	// 		fmt.Println("Not adding because the follower already exists")
	// 	}
	// }
	// fmt.Println("These is the array: ", user["followers"])
	// if !followerExists {
	// 	userIdRaw, _ := primitive.ObjectIDFromHex(userId)
	// 	filter := bson.M{"_id": bson.M{"$eq": userIdRaw}}
	// 	update := bson.M{"$set": user}
	// 	_, err := collection.UpdateOne(context.Background(), filter, update)
	// 	if err != nil {
	// 		fmt.Println("Error occurred while attempting to add follower")
	// 		fmt.Println("This is the err: ", err)
	// 	} else {
	// 		fmt.Println("Added follower")
	// 	}
	// }
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
	// if userId == followingId {
	// 	fmt.Println("User can't add himself as following")
	// }
	// user := findUserById(userId)
	// following := findUserById(followingId)
	// if user == nil {
	// 	fmt.Println("The user doesn't exist: ", userId)
	// }
	// if following == nil {
	// 	fmt.Println("The following doesn't exist: ", followingId)
	// }
	// if (user == nil || following == nil) {
	// 	fmt.Println("One of the user doesn't exist")
	// 	return
	// }
	// var followingExists = false
	// if user["following"] == nil {
	// 	user["following"] = [1]string{followingId}
	// } else {
	// 	var existingFollowing bson.A = user["following"].(bson.A)
	// 	for _, currFollowing := range existingFollowing {
	// 		if currFollowing == followingId {
	// 			followingExists = true
	// 		}
	// 	}
	// 	if !followingExists {
	// 		user["following"] = append(existingFollowing, followingId)
	// 		fmt.Println("Following doesn't exist, hence proceeding to add")
	// 	} else {
	// 		fmt.Println("Not adding because the following already exists")
	// 	}
	// }
	// fmt.Println("This is the following array: ", user["following"])
	// if !followingExists {
	// 	userIdRaw, _ := primitive.ObjectIDFromHex(userId)
	// 	filter := bson.M{"_id": bson.M{"$eq": userIdRaw}}
	// 	update := bson.M{"$set": user}
	// 	_, err := collection.UpdateOne(context.Background(), filter, update)
	// 	if err != nil {
	// 		fmt.Println("Error occurred while attempting to add follower")
	// 		fmt.Println("This is the err: ", err)
	// 	} else {
	// 		fmt.Println("Added following")
	// 	}
	// }
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
	distString := r.URL.Query().Get("dist") // this is in metres
	dist, err := strconv.ParseFloat(distString, 64)
	if err != nil {
		fmt.Println("distance parameter for API is not a float value")
		return
	}
	limitString := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitString)
	fmt.Println("this is limit: ", limit)
	if err != nil {
		fmt.Println("limit parameter for API is not a int value")
		return
	}
	user := findUserById(userId)
	if user == nil {
		fmt.Println("User doesn't exist")
	}
	if user["latitude"] == nil || user["longitude"] == nil {
		fmt.Println("User's location information (lat and long) is not available")
		return
	}
	userLat := user["latitude"].(float64)
	userLong := user["longitude"].(float64)
	var results []primitive.M
	var followingList bson.A
	if user["following"] != nil {
		followingList = user["following"].(bson.A)
		for _, currFollowing := range followingList {
			currUser := findUserById(currFollowing.(string))
			if currUser == nil {
				fmt.Println("[Error] could not retrieve this user: ", currFollowing)
			} else {
				if currUser["latitude"] == nil || currUser["longitude"] == nil {
					continue
				}
				var lat = currUser["latitude"].(float64)
				var long = currUser["longitude"].(float64)
				curr_dist := getDist(userLat, userLong, lat, long)
				if curr_dist < dist {
					results = append(results, currUser)
				}
			}
		}
	} else {
		fmt.Println("This user is not following anyone")
	}
	json.NewEncoder(w).Encode(results)
	fmt.Println("This is the list of nearby friends: ", results)
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
