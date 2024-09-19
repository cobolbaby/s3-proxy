package service

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/xeipuuv/gojsonschema"
)

var (
	jsonSchemaCache = make(map[string]*gojsonschema.Schema)
	cacheLock       sync.RWMutex
)

// Service 是处理请求的主要服务
type Service struct {
	db     *sql.DB
	mc     *minio.Client
	bucket string
}

// NewService 创建并返回 Service 实例
func NewService(db *sql.DB, mc *minio.Client, bucket string) *Service {
	return &Service{
		db:     db,
		mc:     mc,
		bucket: bucket,
	}
}

// HandleRequest 处理上传请求
func (s *Service) HandleRequest(c *gin.Context) {
	// 从 HTTP Header 获取 JSON Schema ID
	schemaID := c.GetHeader("X-Json-Schema-ID")
	if schemaID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "缺少 JSON Schema ID"})
		return
	}

	// 校验 JSON 数据
	if err := s.validateJSON(c, schemaID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 保存 JSON 到 MinIO
	if err := s.saveToMinIO(c, schemaID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "数据上传成功"})
}

// validateJSON 根据 Schema 校验传入的 JSON 数据
func (s *Service) validateJSON(c *gin.Context, schemaID string) error {
	cacheLock.RLock()
	schema, exists := jsonSchemaCache[schemaID]
	cacheLock.RUnlock()

	if !exists {
		// 从数据库加载 JSON Schema
		row := s.db.QueryRow("SELECT schema FROM mc_json_schema WHERE schema_id = $1", schemaID)
		var schemaStr string
		if err := row.Scan(&schemaStr); err != nil {
			return err
		}

		// 编译 Schema 并缓存
		compiledSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaStr))
		if err != nil {
			return err
		}

		cacheLock.Lock()
		jsonSchemaCache[schemaID] = compiledSchema
		cacheLock.Unlock()
		schema = compiledSchema
	}

	// 校验 JSON 数据
	body, err := c.GetRawData()
	if err != nil {
		return err
	}
	documentLoader := gojsonschema.NewBytesLoader(body)
	result, err := schema.Validate(documentLoader)
	if err != nil {
		return err
	}

	if !result.Valid() {
		return errors.New(result.Errors()[0].Description())
	}

	return nil
}

// saveToMinIO 将数据保存到 MinIO
func (s *Service) saveToMinIO(c *gin.Context, schemaID string) error {
	t := time.Now()
	objectName := filepath.Join(t.Format("2006-01-02"), schemaID+".json")

	_, err := s.mc.PutObject(
		context.Background(),
		s.bucket,
		objectName,
		c.Request.Body,
		c.Request.ContentLength,
		minio.PutObjectOptions{ContentType: "application/json"},
	)
	return err
}
