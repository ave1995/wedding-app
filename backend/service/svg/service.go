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

func NewSvgService(store store.SvgStore) service.SvgService {
	return &svgService{store: store}
}

// GetUserSvgs implements service.SvgService.
func (s *svgService) GetUserSvgs(ctx context.Context) ([]*model.SVG, error) {
	return s.store.GetUserSvgs(ctx)
}
