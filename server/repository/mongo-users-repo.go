package repository

import (
	"time"
	"context"
	"../models"
	"log"
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
	dbName = "rydedb"
	collectionName = "users"
)

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



