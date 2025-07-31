package mongodb

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type userStore struct {
	database *mongo.Database
}

func NewUserStore(database *mongo.Database) store.UserStore {
	return &userStore{database: database}
}

func (r *userStore) usersCollection() *mongo.Collection {
	return r.database.Collection("users")
}

func (r *userStore) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	collection := r.usersCollection()

	//TODO: password hash and add

	mongoUser := &user{
		ID:          uuid.New(),
		Username:    username,
		Email:       email,
		IsTemporary: false,
	}

	_, err := collection.InsertOne(ctx, mongoUser)
	if err != nil {
		return nil, err
	}

	return mongoUser.ToDomain(), nil
}
