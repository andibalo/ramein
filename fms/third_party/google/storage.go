package google

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"os"
)

func NewGoogleCloudStorageClient(log *log.Helper) *storage.Client {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "../../ramein-381009-6aa75cc8440d.json")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create google storage client: %v", err)
	}

	return client
}
