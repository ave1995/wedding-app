package googlecloud

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"wedding-app/config"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type svgStore struct {
	store      *storage.Client
	bucketName string
	logger     *slog.Logger
}

func NewCloud(store *storage.Client, config config.BucketConfig, logger *slog.Logger) store.SvgStore {
	return &svgStore{store: store, bucketName: config.UserIconsBucket, logger: logger}
}

// GetUserSvgs implements store.SvgStore.
func (s *svgStore) GetUserSvgs(ctx context.Context) ([]*model.SVG, error) {
	items := s.store.Bucket(s.bucketName).Objects(ctx, nil)

	var svgs []*model.SVG
	for {
		item, err := items.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate to next one: %w", err)
		}
		if strings.HasSuffix(item.Name, ".svg") {
			url := "https://storage.googleapis.com/" + s.bucketName + "/" + item.Name
			nameWithoutSuffix := strings.TrimSuffix(item.Name, ".svg")

			svg := &model.SVG{
				Name: nameWithoutSuffix,
				URL:  url,
			}
			svgs = append(svgs, svg)
		}
	}

	return svgs, nil
}
