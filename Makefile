# 项目名称
APP_NAME = schema-proxy
VERSION = 1.0.0
BUILD_DIR = bin

# Go 参数
GO_CMD = go
GO_BUILD = $(GO_CMD) build
GO_RUN = $(GO_CMD) run
GO_TEST = $(GO_CMD) test
GO_FMT = $(GO_CMD) fmt
GO_GET = $(GO_CMD) get
GO_MOD_TIDY = $(GO_CMD) mod tidy

# 二进制输出文件名
BINARY = $(BUILD_DIR)/$(APP_NAME)

# 构建环境变量
LDFLAGS = -ldflags "-X main.Version=$(VERSION)"

# 数据库迁移路径
MIGRATE_CMD = migrate
MIGRATE_SOURCE = file://db/migrations
DB_URL = postgres://user:password@localhost:5432/yourdb?sslmode=disable

# 默认任务
all: build

# 安装依赖项
deps:
	$(GO_GET) ./...
	$(GO_MOD_TIDY)

# 代码格式化
fmt:
	$(GO_FMT) ./...

# 编译项目
build: fmt
	$(GO_BUILD) $(LDFLAGS) -o $(BINARY) ./cmd/yourapp

# 运行项目
run: build
	$(GO_RUN) ./cmd/yourapp/main.go

# 运行测试
test:
	$(GO_TEST) ./...

# 清理生成的文件
clean:
	rm -rf $(BUILD_DIR)

# 执行数据库迁移
migrate-up:
	$(MIGRATE_CMD) -source $(MIGRATE_SOURCE) -database $(DB_URL) up

# 回滚数据库迁移
migrate-down:
	$(MIGRATE_CMD) -source $(MIGRATE_SOURCE) -database $(DB_URL) down

# 创建新的迁移文件
migrate-create:
	$(MIGRATE_CMD) create -ext sql -dir db/migrations -seq $(name)

# 帮助文档
help:
	@echo "可用命令:"
	@echo "  make all           - 默认任务 (构建项目)"
	@echo "  make deps          - 安装依赖项"
	@echo "  make fmt           - 格式化代码"
	@echo "  make build         - 编译项目"
	@echo "  make run           - 运行项目"
	@echo "  make test          - 运行测试"
	@echo "  make clean         - 清理生成的二进制文件"
	@echo "  make migrate-up    - 执行数据库迁移"
	@echo "  make migrate-down  - 回滚数据库迁移"
	@echo "  make migrate-create name=<migration_name> - 创建新的迁移文件"
