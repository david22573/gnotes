package store

import (
	"context"
	"fmt"
	"time"

	"github.com/david22573/gnotes/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoStore struct {
	client *mongo.Client
}

func NewMongoStore(uri string) (*MongoStore, error) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mongodb: %w", err)
	}

	return &MongoStore{
		client: client,
	}, nil
}

func (s *MongoStore) Get(key string) (*types.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("gnotes").Collection("user")
	filter := bson.M{"_id": key}

	var user types.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return &types.User{}, fmt.Errorf("failed to find user: %w", err)
	}

	return &user, nil
}

func (s *MongoStore) Set(key, user *types.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("gnotes").Collection("user")
	filter := bson.M{"_id": key}

	update := bson.M{"$set": bson.M{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (s *MongoStore) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collection := s.client.Database("gnotes").Collection("users")
	filter := bson.M{"_id": key}

	_, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}
