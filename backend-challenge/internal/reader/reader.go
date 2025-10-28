package reader

import (
	"context"
)

type FileReader interface {
	SearchPromo(ctx context.Context, promo string) (bool, error)
}
