package mongodb

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectClient(ctx context.Context, logger *slog.Logger, url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)

	//Only 10 sec for an attempt
	connectionContext, connectionContextCancel := context.WithTimeout(ctx, 10*time.Second)
	defer connectionContextCancel()

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
