package handlers

import (
	"context"

	"getitqec.com/server/file/pkg/logger"
	"getitqec.com/server/file/pkg/model"
)

type UploadFileHandler struct {
	Model model.FileModelI
}

func (s *UploadFileHandler) UploadFile(ctx context.Context, name string, data []byte) error {
	// userID, err := commons.GetUserID(ctx)
	// if err != nil {
	// 	return nil, err
	// }
	// return s.Model.UploadImage(ctx, name, data)
	if err := s.Model.UploadFile(ctx, name, data); err != nil {
		logger.Log.Error(err.Error())
		return err
	}
	return nil
}
