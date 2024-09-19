package minio

import (
	"context"
	"log"
	"time"

	"github.com/cobolbaby/s3-proxy/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// NewMinioClient 初始化 MinIO 客户端
func NewMinioClient(cfg *config.Config) (*minio.Client, error) {
	// 初始化 MinIO 客户端
	mc, err := minio.New(cfg.Minio.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.Minio.AccessKeyID, cfg.Minio.SecretAccessKey, ""),
		Secure: cfg.Minio.UseSSL,
	})
	if err != nil {
		return nil, err
	}

	// 测试连接和重试机制
	for i := 0; i < 3; i++ {
		_, err = mc.ListBuckets(context.Background())
		if err == nil {
			log.Println("MinIO 客户端初始化成功")
			return mc, nil
		}
		log.Printf("MinIO 客户端连接失败，重试中 (%d/3)...", i+1)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
