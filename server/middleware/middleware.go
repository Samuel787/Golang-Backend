package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"../models"
	"github.com/gorilla/mux"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://SamuelRyde:SlfZ0ehN2bDKnr3h@rydecluster.pbok3.mongodb.net/test?retryWrites=true&w=majority"
const dbName = "rydedb"
const collectionName = "users"
var collection *mongo.Collection

func init() {
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOptions) // mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	
	database := client.Database(dbName)
	collection = database.Collection(collectionName)
	
}

// Create User route
func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	insertOneUser(user)
	json.NewEncoder(w).Encode(user)
}

func insertOneUser(user models.User) {
	insertResult, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted a Single Record ", insertResult.InsertedID)
}

// Get All Users route
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	payload := getAllUsers()
	json.NewEncoder(w).Encode(payload)
}

// Get User by their ID
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)
	payload := findUserById(params["id"])
	json.NewEncoder(w).Encode(payload)
}

func findUserById(user string) bson.M {
	var result bson.M
	id, _ := primitive.ObjectIDFromHex(user)
	collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	return result;
}

// get all users from the DB and return it
func getAllUsers() []primitive.M {
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var results []primitive.M
	for cur.Next(context.Background()) {
		var result bson.M
		e := cur.Decode(&result)
		if e != nil {
			log.Fatal(e)
		}
		results = append(results, result)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	cur.Close(context.Background())
	return results
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	params := mux.Vars(r)
	deleteOneUser(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func deleteOneUser(user string) {
	fmt.Println(user)
	id, _ := primitive.ObjectIDFromHex(user)
	filter := bson.M{"_id": id}
	d, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Document", d.DeletedCount)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "PUT")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	// params := mux.Vars(r)
	hex_id := r.URL.Query().Get("id")
	if hex_id != "" {
		updateOneUser(r)
	}

}

func updateOneUser(r *http.Request) {
	hex_id := r.URL.Query().Get("id")
	id, _ := primitive.ObjectIDFromHex(hex_id)
	filter := bson.M{"_id": bson.M{"$eq": id}}
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
	update := bson.M{"$set" : change}
	_, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		fmt.Println("UpdateOneUser() result ERROR: ", err)
	} else {
		fmt.Println("Updated record successfully")
	}
}
