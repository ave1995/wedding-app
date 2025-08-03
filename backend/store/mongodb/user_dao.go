package mongodb

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/store"
	"wedding-app/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type userStore struct {
	database *mongo.Database
}

func NewUserStore(database *mongo.Database) store.UserStore {
	return &userStore{database: database}
}

const Users_Collection = "users"

func (r *userStore) usersCollection() *mongo.Collection {
	return r.database.Collection(Users_Collection)
}

func (r *userStore) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	collection := r.usersCollection()

	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	mongoUser := &user{
		ID:          uuid.New(),
		Username:    username,
		Email:       email,
		Password:    hashedPass,
		IsTemporary: false,
	}

	_, err = collection.InsertOne(ctx, mongoUser)
	if err != nil {
		return nil, err
	}

	return mongoUser.ToDomain(), nil
}
