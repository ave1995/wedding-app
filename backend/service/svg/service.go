package svg

import (
	"context"
	"wedding-app/domain/model"
	"wedding-app/domain/service"
	"wedding-app/domain/store"
)

type svgService struct {
	store store.SvgStore
}

func NewSvgService(store store.SvgStore) service.StorageService {
	return &svgService{store: store}
}

// GetUserSvgs implements service.SvgService.
func (s *svgService) GetUserSvgs(ctx context.Context) ([]*model.BucketItemUrl, error) {
	return s.store.GetUserSvgs(ctx)
}

// GetBucketData implements service.StorageService.
func (s *svgService) GetBucketData(ctx context.Context, bucketName string, bucketItemName string) (*model.BucketItemData, error) {
	return s.store.GetBucketData(ctx, bucketName, bucketItemName)
}

// GetBucketUrls implements service.StorageService.
func (s *svgService) GetBucketUrls(ctx context.Context, bucketName string, suffix string) ([]*model.BucketItemUrl, error) {
	return s.store.GetBucketUrls(ctx, bucketName, suffix)
}
