package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"wedding-app/domain/model"
	"wedding-app/domain/store"
	"wedding-app/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
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

func (s *userStore) usersCollection() *mongo.Collection {
	return s.database.Collection(UsersCollection)
}

func (s *userStore) RegisterUser(ctx context.Context, params model.RegisterUserParams) (*model.User, error) {
	collection := s.usersCollection()

	hashedPass, err := utils.HashPassword(params.Password)
	if err != nil {
		return nil, err
	}

	mongoUser := &user{
		ID:          uuid.NewString(),
		Username:    params.Username,
		Email:       params.Email,
		Password:    hashedPass,
		IsTemporary: false,
		IconUrl:     params.IconURL,
	}

	_, err = collection.InsertOne(ctx, mongoUser)
	if err != nil {
		s.logger.Error("failed to insert one to: %w", utils.ErrAttr(err), slog.Any("user", mongoUser))
		return nil, err
	}

	user, err := mongoUser.ToDomain()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// LoginUser implements store.UserStore.
func (s *userStore) LoginUser(ctx context.Context, email string, password string) (*model.User, error) {
	collection := s.usersCollection()

	var mongoUser user
	err := collection.FindOne(ctx, bson.M{"email": email}).Decode(&mongoUser)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// User not found
			return nil, fmt.Errorf("invalid email")
		}
		return nil, err
	}

	if !utils.CheckPasswordHash(password, mongoUser.Password) {
		return nil, fmt.Errorf("invalid password")
	}

	user, err := mongoUser.ToDomain()
	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetUserByID implements store.UserStore.
func (s *userStore) GetUserByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	result, err := getByFilter[user](ctx, s.usersCollection(), bson.M{userFieldID: id.String()})
	if err != nil {
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	user, err := result.ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert user to domain model: %w", err)
	}

	return user, nil
}
