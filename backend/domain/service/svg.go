package service

import (
	"context"
	"wedding-app/domain/model"
)

type SvgService interface {
	GetUserSvgs(ctx context.Context) ([]*model.SVG, error)
}
