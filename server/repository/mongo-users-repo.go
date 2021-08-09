package repository

import (
	"context"
	"errors"
	"log"
	"time"

	"../models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type repo struct{}

// NewMongoUsersRepository creates a new repo
func NewMongoUsersRepository() UsersRepository {
	return &repo{}
}

const (
	connectionString = "mongodb+srv://SamuelRyde:SlfZ0ehN2bDKnr3h@rydecluster.pbok3.mongodb.net/test?retryWrites=true&w=majority"
	dbName           = "rydedb"
	collectionName   = "users"
)

func ConnectToDatabase() (mongo.Client, context.Context, error, mongo.Collection) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	return *client, ctx, err, *collection
}

func (*repo) AddUser(user *models.User) (*models.User, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return user, nil
	}
}

func (*repo) GetAllUsers() ([]primitive.M, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	database := client.Database(dbName)
	collection := database.Collection(collectionName)
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
	if err != nil {
		return nil, err
	} else {
		return results, nil
	}
}

func (*repo) GetUser(userId string) (bson.M, error) {
	client, ctx, err, collection := ConnectToDatabase()
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var result bson.M
	id, _ := primitive.ObjectIDFromHex(userId)
	collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	if result == nil {
		return nil, errors.New("[Get User] User does not exist in the database")
	} else {
		return result, nil
	}
}

func (*repo) DeleteUser(userId string) error {
	client, ctx, err, collection := ConnectToDatabase()
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	id, _ := primitive.ObjectIDFromHex(userId)
	var resultUser bson.M
	filter := bson.M{"_id": id}
	collection.FindOne(context.Background(), filter).Decode(&resultUser)
	if resultUser == nil {
		return errors.New("[DeleteUser] user to delete does not exist in database")
	} else {
		_, resultErr := collection.DeleteOne(context.Background(), filter)
		return resultErr
	}
}

func (*repo) UpdateUser(userId string, update bson.M) error {
	client, ctx, err, collection := ConnectToDatabase()
	defer client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
		return err
	}
	id, _ := primitive.ObjectIDFromHex(userId)
	var resultUser bson.M
	filter := bson.M{"_id": id}
	collection.FindOne(context.Background(), filter).Decode(&resultUser)
	if resultUser == nil {
		return errors.New("[UpdateUser] user to update does not exist in database")
	} else {
		userFilter := bson.M{"_id": bson.M{"$eq": id}}
		_, err := collection.UpdateOne(context.Background(), userFilter, update)
		return err
	}
}
