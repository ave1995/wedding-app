package mongodb

import (
	"context"
	"log/slog"
	"wedding-app/domain/model"
	"wedding-app/domain/store"
	"wedding-app/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type userStore struct {
	database *mongo.Database
	logger   *slog.Logger
}

func NewUserStore(database *mongo.Database, logger *slog.Logger) store.UserStore {
	return &userStore{database: database, logger: logger}
}

const UsersCollection = "users"

func (r *userStore) usersCollection() *mongo.Collection {
	return r.database.Collection(UsersCollection)
}

func (r *userStore) RegisterUser(ctx context.Context, username, email, password string) (*model.User, error) {
	collection := r.usersCollection()

	hashedPass, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	mongoUser := &user{
		ID:          uuid.New().String(),
		Username:    username,
		Email:       email,
		Password:    hashedPass,
		IsTemporary: false,
	}

	_, err = collection.InsertOne(ctx, mongoUser)
	if err != nil {
		r.logger.Error("failed to insert one to: %w", utils.ErrAttr(err), slog.Any("user", mongoUser))
		return nil, err
	}

	user, err := mongoUser.ToDomain()
	if err != nil {
		return nil, err
	}

	return user, nil
}
