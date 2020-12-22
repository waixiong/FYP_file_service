package handlers

import (
	"context"

	"getitqec.com/server/file/pkg/logger"
	"getitqec.com/server/file/pkg/model"
)

type UploadImageHandler struct {
	Model model.FileModelI
}

func (s *UploadImageHandler) UploadImage(ctx context.Context, bucketName string, name string, data []byte) (string, error) {
	// userID, err := commons.GetUserID(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// return s.Model.UploadImage(ctx, name, data)
	id, err := s.Model.UploadImage(ctx, bucketName, name, data)
	if err != nil {
		logger.Log.Error(err.Error())
		return "", err
	}
	return id, nil
}
