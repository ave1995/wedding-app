package mongodb

import (
	"context"
	"wedding-app/models/domain"
	"wedding-app/models/mongodb"
	"wedding-app/repositories"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type userMongoRepostiory struct {
	database *mongo.Database
}

func NewUserRepository(database *mongo.Database) repositories.UserRepository {
	return &userMongoRepostiory{database: database}
}

func (r *userMongoRepostiory) usersCollection() *mongo.Collection {
	return r.database.Collection("users")
}

func (r *userMongoRepostiory) RegisterUser(ctx context.Context, username, email, password string) (*domain.User, error) {
	collection := r.usersCollection()

	//TODO: password hash and add

	mongoUser := mongodb.MongoUser{
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
