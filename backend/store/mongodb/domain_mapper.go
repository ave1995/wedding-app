package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DomainMapper[D any] interface {
	ToDomain() (*D, error)
}

// getByFilterAndConvert fetches a single document from a MongoDB collection
// that matches the given filter, decodes it into a local model struct, and converts
// it to the corresponding domain model struct.
//
// Type Parameters:
//
//	L - The local model type that represents the MongoDB document structure.
//	    Must implement DomainMapper[D] so it can be converted to the domain model.
//	D - The domain model type that represents the application's core business entity.
//
// Parameters:
//
//	ctx        - The context for managing request deadlines and cancellation.
//	collection - The MongoDB collection to query.
//	filter     - The BSON filter used to select the document.
//
// Returns:
//
//	A pointer to the domain model (D) or an error if fetching, decoding, or converting fails.
func getByFilterAndConvert[L DomainMapper[D], D any](
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
) (*D, error) {
	result, err := getByFilter[L](ctx, collection, filter)
	if err != nil {
		return nil, err
	}

	domain, err := (*result).ToDomain()
	if err != nil {
		return nil, fmt.Errorf("failed to convert to domain: %w", err)
	}

	return domain, nil
}

// getManyByFilterAndConvert fetches multiple documents from a MongoDB collection
// that match the given filter, decodes them into local model structs, and converts
// them to domain model structs.
//
// Type Parameters:
//
//	L - The local model type that represents the MongoDB document structure.
//	    Must implement DomainMapper[D] so it can be converted to the domain model.
//	D - The domain model type that represents the application's core business entity.
//
// Parameters:
//
//	ctx        - The context for managing request deadlines and cancellation.
//	collection - The MongoDB collection to query.
//	filter     - The BSON filter used to select the document.
//
// Returns:
//
//	A slice of pointers to the domain model (D) or an error if fetching, decoding, or converting fails.
func getManyByFilterAndConvert[L DomainMapper[D], D any](
	ctx context.Context,
	collection *mongo.Collection,
	filter bson.M,
	sort *bson.D,
) ([]*D, error) {
	var opts *options.FindOptions
	if sort != nil {
		opts = options.Find().SetSort(sort)
	}
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get cursor: %w", err)
	}
	defer cursor.Close(ctx)

	var results []L
	if err = cursor.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("failed to get all results: %w", err)
	}

	domains := make([]*D, len(results))
	for i, r := range results {
		domain, err := r.ToDomain()
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain: %w", err)
		}
		domains[i] = domain
	}

	return domains, nil
}

// createAndConvert inserts a MongoDB document represented by a local model into the given collection
// and converts it into the corresponding domain model.
//
// Type Parameters:
//
//	L - The local MongoDB model type. Must implement DomainMapper[D] to allow conversion to the domain model.
//	D - The domain model type representing the application's core business entity.
//
// Parameters:
//
//	ctx         - Context for request scoping, deadlines, and cancellation.
//	collection  - The MongoDB collection where the document will be inserted.
//	mongoModel  - The local model instance to insert into MongoDB.
//
// Returns:
//
//	A pointer to the domain model (D) corresponding to the inserted document,
//	or an error if insertion or conversion fails.
func createAndConvert[L DomainMapper[D], D any](
	ctx context.Context,
	collection *mongo.Collection,
	mongoModel L,
) (*D, error) {
	_, err := collection.InsertOne(ctx, mongoModel)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	domain, err := mongoModel.ToDomain()
	if err != nil {
		return nil, err
	}

	return domain, nil
}
