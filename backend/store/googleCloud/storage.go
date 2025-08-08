package googlecloud

import (
	"context"

	"cloud.google.com/go/storage"
)

func ConnectClient(ctx context.Context) (*storage.Client, error) {
	return storage.NewClient(ctx)
}
