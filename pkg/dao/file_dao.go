package dao

import (
	"context"

	"getitqec.com/server/file/pkg/commons"
)

type FileDAO struct {
	mongodb commons.MongoDB
}

func InitFileDAO(m commons.MongoDB) IFileDAO {
	return &FileDAO{mongodb: m}
}

func (v *FileDAO) Create(ctx context.Context, name string, data []byte) error {
	return v.mongodb.Upload(ctx, name, data)
}

func (v *FileDAO) Get(ctx context.Context, name string) ([]byte, error) {
	return v.mongodb.Download(ctx, name)
}

func (v *FileDAO) Delete(ctx context.Context, name string) error {
	return nil
}
