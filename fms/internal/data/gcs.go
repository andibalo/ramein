package data

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"io"
)

type GCSRepo interface {
	Upload(file []byte, bucketName, fileName, filePath string) error
}

type googleCloudStorageRepo struct {
	gcs *storage.Client
	log *log.Helper
}

func NewGoogleCloudStorageRepo(gcs *storage.Client, logger log.Logger) *googleCloudStorageRepo {
	return &googleCloudStorageRepo{
		gcs: gcs,
		log: log.NewHelper(logger),
	}
}

func (r *googleCloudStorageRepo) Upload(file []byte, bucketName, fileName, filePath string) error {

	wc := r.gcs.Bucket(bucketName).Object(filePath + fileName).NewWriter(context.Background())
	if _, err := io.Copy(wc, bytes.NewBuffer(file)); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
