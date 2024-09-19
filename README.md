
# YourProject

YourProject 是一个基于 Go 的 Web 服务，用户可以通过 API 上传 JSON 数据，服务会根据 PostgreSQL 中存储的 JSON Schema 校验数据，并将通过校验的数据保存到 MinIO 存储中。

## 功能概述

1. 接收 JSON 数据，通过 HTTP POST 请求上传
2. 从 PostgreSQL 数据库中获取并缓存 JSON Schema，根据 Schema 校验上传的数据
3. 校验通过后，将 JSON 数据保存到 MinIO 的指定 bucket 中，每天创建一个新的子目录
4. 支持通过 Kubernetes ConfigMap 动态配置数据库和 MinIO 连接参数
5. 支持数据库迁移，类似 Flyway

## 项目结构

```text
yourproject/
│
├── cmd/
│   └── yourapp/
│       └── main.go          # 应用的入口文件
│
├── internal/
│   ├── database/
│   │   └── database.go      # 数据库相关逻辑
│   ├── minio/
│   │   └── minio_client.go  # MinIO 相关逻辑
│   └── service/
│       └── service.go       # 主要业务逻辑
│
├── pkg/
│   ├── config/
│   │   └── config.go        # 配置加载逻辑
│   └── migration/
│       └── migration.go     # 数据库迁移逻辑
│
├── db/
│   └── migrations/          # 数据库迁移文件（SQL）
│       └── 001_create_schema.sql
│
└── go.mod                   # Go 模块管理文件
└── go.sum                   # 依赖的校验文件
```

## 配置文件

项目使用 [Viper](https://github.com/spf13/viper) 加载配置，支持通过环境变量或 Kubernetes ConfigMap 动态配置。示例配置文件如下：

```yaml
database:
  url: "postgres://user:password@localhost:5432/yourdb?sslmode=disable"

minio:
  endpoint: "play.min.io"
  accessKeyID: "minioadmin"
  secretAccessKey: "minioadmin"
  useSSL: true
  bucket: "your-bucket-name"
```

## 使用方法

### 1. 安装依赖

首先安装依赖项：

```bash
make deps
```

### 2. 编译项目

```bash
make build
```

### 3. 运行项目

```bash
make run
```

项目会启动一个 HTTP 服务，监听在 `:8080` 端口。你可以通过 POST 请求上传 JSON 数据：

```bash
curl -X POST -H "X-Json-Schema-ID: schema123" -H "Content-Type: application/json" -d '{"key": "value"}' http://localhost:8080/upload
```

### 4. 运行测试

```bash
make test
```

### 5. 数据库迁移

#### 创建迁移文件

要创建一个新的迁移文件，可以运行以下命令：

```bash
make migrate-create name=create_new_table
```

#### 执行迁移

```bash
make migrate-up
```

#### 回滚迁移

```bash
make migrate-down
```

## 部署

可以通过 Kubernetes ConfigMap 动态注入配置，确保 `config.yaml` 中的内容符合你的部署环境。

### 环境变量

你可以通过以下环境变量动态配置服务：

- `DATABASE_URL`：PostgreSQL 数据库连接 URL
- `MINIO_ENDPOINT`：MinIO 服务的地址
- `MINIO_ACCESS_KEY_ID`：MinIO 的访问密钥
- `MINIO_SECRET_ACCESS_KEY`：MinIO 的访问密钥
- `MINIO_USE_SSL`：是否使用 SSL 连接 MinIO
- `MINIO_BUCKET`：存储 JSON 数据的 bucket 名称

## 贡献

欢迎提出 issues 或者提交 pull request 来改进该项目。

## 许可证

本项目使用 [MIT License](LICENSE)。
