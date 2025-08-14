package mongodb

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wedding-app/config"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectClient(ctx context.Context, logger *slog.Logger, config config.StoreConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.DbUrl)

	if config.DbUsername != "" && config.DbPassword != "" {
		clientOptions.SetAuth(options.Credential{
			Username: config.DbUsername,
			Password: config.DbPassword,
		})
	}

	//Only 10 sec for an attempt
	connectionContext, connectionContextCancel := context.WithTimeout(ctx, 10*time.Second)
	defer connectionContextCancel()

	logger.Info("Attempting to connect to MongoDB", slog.String("url", config.DbUrl))

	client, err := mongo.Connect(connectionContext, clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	logger.Info("MongoClient connected")

	return client, nil
}

func GetDatabase(client *mongo.Client, name string) *mongo.Database {
	return client.Database(name)
}

// getByFilter fetches a single document from a MongoDB collection that matches the given filter
// and decodes it into the specified local model type.
//
// Type Parameters:
//
//	L - The local model type that represents the MongoDB document structure.
//
// Parameters:
//
//	ctx        - The context for managing request deadlines and cancellation.
//	collection - The MongoDB collection to query.
//	filter     - The BSON filter used to select the document.
//
// Returns:
//
//	A pointer to a localModel instance containing the decoded document,
//	or an error if the query fails or no document is found.
func getByFilter[L any](ctx context.Context, collection *mongo.Collection, filter bson.M) (*L, error) {
	var result L
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, mongo.ErrNoDocuments
		}
		return nil, fmt.Errorf("find by filter failed: %w", err)
	}
	return &result, nil
}

const FieldID = "_id"
