package infra

import (
	"context"
	"log"
	"time"

	"github.com/BernardoDenkvitts/MongoAPP/internal/types"
	"github.com/BernardoDenkvitts/MongoAPP/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (s *MongoDBStore) Init() {

	collections, err := s.Db.ListCollectionNames(context.TODO(), bson.M{"name": "user"})
	utils.FailOnError(err, "Fail to get collections")

	if len(collections) == 0 {
		validator := bson.M{
			"$jsonSchema": userSchema(),
		}
		opts := options.CreateCollection().SetValidator(validator)

		err = s.Db.CreateCollection(context.TODO(), "user", opts)
		utils.FailOnError(err, "Error to initiate database")

		createIndex(s.Db)
	}

	log.Println("Database Initialized")
}

func (s *MongoDBStore) getUserCollection() *mongo.Collection {
	return s.Db.Collection(collectionName)
}

func (s *MongoDBStore) CreateUserInformation(user *types.User) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	_, err := s.getUserCollection().InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoDBStore) GetUserById(id string) (*types.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	findUser := s.getUserCollection().FindOne(ctx, bson.M{"id": id})
	if err := findUser.Err(); err != nil {
		return nil, nil
	}

	user := new(types.User)
	err := findUser.Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *MongoDBStore) GetUsersInformations() ([]*types.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	cursor, err := s.getUserCollection().Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	users, err := decodeUsersFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *MongoDBStore) GetLatestUserInformations() ([]*types.User, error) {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	fiveMinutesAgo := time.Now().Add(-5 * time.Minute)

	filter := bson.M{
		"created_at": bson.M{"$gt": fiveMinutesAgo},
	}

	cursor, err := s.getUserCollection().Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	latestUsers, err := decodeUsersFromCursor(ctx, cursor)
	if err != nil {
		return nil, err
	}

	return latestUsers, nil
}

func decodeUsersFromCursor(ctx context.Context, cursor *mongo.Cursor) ([]*types.User, error) {
	var users []*types.User
	for cursor.Next(ctx) {
		user := new(types.User)
		err := cursor.Decode(&user)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
