package store

import (
	"context"
	"wedding-app/domain/model"
)

type SvgStore interface {
	GetUserSvgs(ctx context.Context) ([]*model.BucketItemUrl, error)
	GetBucketUrls(ctx context.Context, bucketName string, suffix string) ([]*model.BucketItemUrl, error)
	GetBucketData(ctx context.Context, bucketName string, bucketItemName string) (*model.BucketItemData, error)
}
