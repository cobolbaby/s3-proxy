package main

import (
	"log"

	"github.com/cobolbaby/schema-proxy/internal/database"
	"github.com/cobolbaby/schema-proxy/internal/minio"
	"github.com/cobolbaby/schema-proxy/internal/service"
	"github.com/cobolbaby/schema-proxy/pkg/config"
	"github.com/cobolbaby/schema-proxy/pkg/migration"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置文件
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置文件失败: %v", err)
	}

	// 执行数据库迁移
	migration.RunMigrations(cfg.Database.URL)

	// 初始化数据库和 MinIO 客户端
	db, err := database.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	mc, err := minio.NewMinioClient(cfg)
	if err != nil {
		log.Fatalf("MinIO 客户端初始化失败: %v", err)
	}

	// 创建 Service 实例并注入依赖
	svc := service.NewService(db, mc, cfg.Minio.Bucket)

	// 启动 Gin web 服务
	r := gin.Default()
	r.POST("/upload", svc.HandleRequest)
	r.Run(":8080")
}
