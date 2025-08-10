package store

import (
	"context"
	"wedding-app/domain/model"
)

type SvgStore interface {
	GetUserSvgs(ctx context.Context) ([]*model.SVG, error)
}
