package config

import (
	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Database DatabaseConfig
	Minio    MinioConfig
}

// DatabaseConfig 数据库配置结构体
type DatabaseConfig struct {
	URL string
}

// MinioConfig MinIO 配置结构体
type MinioConfig struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	UseSSL          bool
	Bucket          string
}

// LoadConfig 加载配置文件
func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
