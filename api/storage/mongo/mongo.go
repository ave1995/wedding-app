package mongo

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectClient(url string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(url)

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	}

	slog.Info("MongoClient connected")

	return client, nil
}

func GetDatabase(client *mongo.Client, name string) *mongo.Database {
	return client.Database(name)
}
