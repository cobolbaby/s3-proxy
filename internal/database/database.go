package database

import (
	"database/sql"
	"log"

	"github.com/cobolbaby/schema-proxy/pkg/config"
	_ "github.com/lib/pq"
)

// NewDatabase 初始化数据库连接
func NewDatabase(cfg *config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.Database.URL)
	if err != nil {
		return nil, err
	}

	// 测试数据库连接
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("数据库连接成功")
	return db, nil
}
