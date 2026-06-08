package storage

import (
	"context"

	"github.com/banggibima/go-english-test-platform/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinIO struct {
	Client *minio.Client
	Bucket string
}

func NewMinIO(cfg *config.Config) (*MinIO, error) {
	client, err := minio.New(
		cfg.MinIOEndpoint,
		&minio.Options{
			Creds: credentials.NewStaticV4(
				cfg.MinIOAccessKey,
				cfg.MinIOSecretKey,
				"",
			),
			Secure: cfg.MinIOUseSSL,
		},
	)

	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(
		context.Background(),
		cfg.MinIOBucket,
	)

	if err != nil {
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(
			context.Background(),
			cfg.MinIOBucket,
			minio.MakeBucketOptions{},
		)

		if err != nil {
			return nil, err
		}
	}

	return &MinIO{
		Client: client,
		Bucket: cfg.MinIOBucket,
	}, nil
}
