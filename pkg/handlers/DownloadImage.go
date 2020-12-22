package handlers

import (
	"context"

	"getitqec.com/server/file/pkg/logger"
	"getitqec.com/server/file/pkg/model"
)

type DownloadImageHandler struct {
	Model model.FileModelI
}

func (s *DownloadImageHandler) DownloadImage(ctx context.Context, bucketName string, name string) ([]byte, error) {
	// userID, err := commons.GetUserID(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	bytes, err := s.Model.DownloadFile(ctx, name)
	if err != nil {
		logger.Log.Error(err.Error())
		return nil, err
	}
	return bytes, nil
}
