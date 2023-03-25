package biz

import (
	"errors"
	v1 "fms/api/file/v1"
	"fms/internal/data"
	"github.com/go-kratos/kratos/v2/log"
)

type FileUsecase struct {
	log     *log.Helper
	gcsRepo data.GCSRepo
}

// NewGreeterUsecase new a Greeter usecase.
func NewFileUsecase(logger log.Logger, gcsRepo data.GCSRepo) *FileUsecase {
	return &FileUsecase{log: log.NewHelper(logger), gcsRepo: gcsRepo}
}

// CreateGreeter creates a Greeter, and returns the new Greeter.
func (uc *FileUsecase) UploadFile(req *v1.UploadFileRequest) error {
	uc.log.Infof("UploadFile: %v", req)

	err := uc.gcsRepo.Upload(req.File, req.Bucket, req.FileName, req.FilePath)

	if err != nil {
		uc.log.Error("Failed to upload file to gcs: %v", err)
		return errors.New("Failed to upload file to gcs")
	}

	return nil
}
