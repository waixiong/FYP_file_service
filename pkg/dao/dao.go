package dao

import (
	"context"
	// "getitqec.com/server/file/pkg/dto"
)

type IFileDAO interface {
	Create(ctx context.Context, name string, data []byte) error
	Get(ctx context.Context, name string) ([]byte, error)
	// Update(ctx context.Context, name string, data) error
	Delete(ctx context.Context, name string) error
}
