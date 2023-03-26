package external

import (
	"bytes"
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/andibalo/ramein/corvus/internal/config"
	"go.uber.org/zap"
	"io"
	"os"
)

type GCSRepo interface {
	Upload(file []byte, bucketName, fileName, filePath string) error
}

type googleCloudStorageRepo struct {
	gcs *storage.Client
	cfg config.Config
}

func NewGoogleCloudStorageClient(cfg config.Config) *storage.Client {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../ramein-381009-6aa75cc8440d.json")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		cfg.Logger().Error("Failed to create google storage client", zap.Error(err))
	}

	return client
}

func NewGoogleCloudStorageRepo(gcs *storage.Client, cfg config.Config) *googleCloudStorageRepo {
	return &googleCloudStorageRepo{
		gcs: gcs,
		cfg: cfg,
	}
}

func (r *googleCloudStorageRepo) Upload(file []byte, bucketName, fileName, filePath string) error {

	wc := r.gcs.Bucket(bucketName).Object(filePath + fileName).NewWriter(context.Background())
	if _, err := io.Copy(wc, bytes.NewBuffer(file)); err != nil {
		r.cfg.Logger().Error(fmt.Sprintf("io.Copy: %v", err), zap.Error(err))
		return fmt.Errorf("io.Copy: %v", err)
	}

	if err := wc.Close(); err != nil {
		r.cfg.Logger().Error(fmt.Sprintf("Writer.Close: %v", err), zap.Error(err))
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}
