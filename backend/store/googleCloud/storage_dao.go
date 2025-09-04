package googlecloud

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"wedding-app/api/restapi"
	"wedding-app/config"
	"wedding-app/domain/model"
	"wedding-app/domain/store"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

type storageStore struct {
	store      *storage.Client
	bucketName string
	logger     *slog.Logger
}

func NewCloud(store *storage.Client, config config.BucketConfig, logger *slog.Logger) store.SvgStore {
	return &storageStore{store: store, bucketName: config.UserIconsBucket, logger: logger}
}

// GetUserSvgs implements store.SvgStore.
func (s *storageStore) GetUserSvgs(ctx context.Context) ([]*model.BucketItemUrl, error) {
	bucket := s.store.Bucket(s.bucketName)
	items := bucket.Objects(ctx, nil)

	var svgs []*model.BucketItemUrl
	for {
		item, err := items.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate to next one: %w", err)
		}
		if strings.HasSuffix(item.Name, ".svg") {
			url := "https://storage.googleapis.com/" + s.bucketName + "/" + item.Name
			nameWithoutSuffix := strings.TrimSuffix(item.Name, ".svg")

			svg := &model.BucketItemUrl{
				Name: nameWithoutSuffix,
				URL:  url,
			}
			svgs = append(svgs, svg)

			rc, err := bucket.Object(item.Name).NewReader(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to create reader: %w", err)
			}

			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("failed to read object: %w", err)
			}

			s.logger.Info("svg file read", slog.String("name", item.Name), slog.Int("lenght", len(data)))

		}
	}

	return svgs, nil
}

// GetBucketUrls implements store.SvgStore.
func (s *storageStore) GetBucketUrls(ctx context.Context, bucketName string, suffix string) ([]*model.BucketItemUrl, error) {
	bucket := s.store.Bucket(s.bucketName)
	items := bucket.Objects(ctx, nil)

	var bucketItems []*model.BucketItemUrl
	for {
		item, err := items.Next()
		if errors.Is(err, iterator.Done) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate to next one: %w", err)
		}
		if strings.HasSuffix(item.Name, suffix) {
			url := restapi.MakeBucketDataURL(bucketName, item.Name)
			bucketItem := &model.BucketItemUrl{
				Name: item.Name,
				URL:  url,
			}

			bucketItems = append(bucketItems, bucketItem)
		}
	}

	return bucketItems, nil
}

// GetBucketData implements store.SvgStore.
func (s *storageStore) GetBucketData(ctx context.Context, bucketName string, bucketItemName string) (*model.BucketItemData, error) {
	bucket := s.store.Bucket(s.bucketName)
	item := bucket.Object(bucketItemName)

	rc, err := item.NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create reader: %w", err)
	}

	data, err := io.ReadAll(rc)
	rc.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	attrs, err := item.Attrs(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read attributes: %w", err)
	}

	return &model.BucketItemData{
		Name:  bucketItemName,
		Data:  data,
		CType: attrs.ContentType,
	}, nil
}
