package infra

import (
	"context"
	"fmt"
	"time"

	"github.com/BernardoDenkvitts/MongoAPP/internal/types"
	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	username       = "root"
	password       = "root"
	hostname       = "localhost:27017"
	dbName         = "mongodbuser"
	collectionName = "user"
)

type Storage interface {
	Init()
	CreateUserInformation(*types.User) error
	GetUserById(id string) (*types.User, error)
	// This function will be use to get users created in the last 5 minutes
	// to be sent to rabbitMQ
	GetLatestUserInformations() ([]*types.User, error)
	GetUsersInformations() ([]*types.User, error)
}

type MongoDBStore struct {
	Client *mongo.Client
	Db     *mongo.Database
}

func dsn() string {
	return fmt.Sprintf("mongodb://%s:%s@%s", username, password, hostname)
}

func NewMongoDBStore() (*MongoDBStore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn()))
	if err != nil {
		return nil, err
	}

	return &MongoDBStore{
		Client: client,
		Db:     client.Database(dbName),
	}, nil
}

func userSchema() primitive.M {
	schema := bson.M{
		"bsonType": "object",
		"required": []string{"id", "firstName", "lastName", "created_at"},
		"properties": bson.M{
			"id": bson.M{
				"bsonType": "string",
			},
			"firstName": bson.M{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z]{1,50}$",
			},
			"lastName": bson.M{
				"bsonType": "string",
				"pattern":  "^[a-zA-Z]{1,50}$",
			},
			"created_at": bson.M{
				"bsonType": "date",
			},
		},
	}

	return schema
}

func createIndex(s *mongo.Database) {
	collection := s.Collection(collectionName)
	index := mongo.IndexModel{
		Keys:    bson.D{{Key: "id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	_, err := collection.Indexes().CreateOne(context.TODO(), index)

	utils.FailOnError(err, "Error to create user index")
}
