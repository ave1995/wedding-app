package mongodb

import (
	"context"
	"log/slog"
	"time"
	"wedding-app/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectClient(ctx context.Context, logger *slog.Logger, config config.StoreConfig) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(config.DbUrl).SetAuth(options.Credential{
		Username: config.DbUsername,
		Password: config.DbPassword,
	})

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
