package model

import (
	"context"
)

type FileModelI interface {
	// sign in ID (email, hp number) to usedID
	UploadFile(ctx context.Context, name string, data []byte) error
	DownloadFile(ctx context.Context, name string) ([]byte, error)
	// UploadPicture()

	UploadImage(ctx context.Context, bucketName string, name string, data []byte) (string, error)
}
