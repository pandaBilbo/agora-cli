# 项目名称和版本
APP_NAME := devex
VERSION := $(shell git describe --tags --always --dirty || echo "dev")
BUILD_TIME := $(shell date +%Y-%m-%d\ %H:%M:%S)
COMMIT_HASH := $(shell git rev-parse --short HEAD || echo "unknown")

# 构建信息
LDFLAGS := -ldflags "-X 'devex/cmd.Version=$(VERSION)' -X 'devex/cmd.BuildTime=$(BUILD_TIME)' -X 'devex/cmd.CommitHash=$(COMMIT_HASH)'"

# 目标平台
PLATFORMS := darwin/amd64 darwin/arm64 linux/amd64 linux/arm64 windows/amd64

# 默认目标
.PHONY: all
all: clean build

# 构建当前平台
.PHONY: build
build:
	@echo "构建 $(APP_NAME) v$(VERSION)..."
	go build $(LDFLAGS) -o bin/$(APP_NAME) main.go

# 构建所有平台
.PHONY: build-all
build-all: clean
	@echo "构建所有平台的 $(APP_NAME) v$(VERSION)..."
	@mkdir -p dist
	@for platform in $(PLATFORMS); do \
		os=$$(echo $$platform | cut -d'/' -f1); \
		arch=$$(echo $$platform | cut -d'/' -f2); \
		output_name=$(APP_NAME); \
		if [ $$os = "windows" ]; then output_name=$${output_name}.exe; fi; \
		echo "构建 $$os/$$arch..."; \
		GOOS=$$os GOARCH=$$arch go build $(LDFLAGS) -o dist/$(APP_NAME)-$$os-$$arch/$$output_name main.go; \
		if [ $$? -eq 0 ]; then \
			cd dist && tar -czf $(APP_NAME)-$$os-$$arch.tar.gz $(APP_NAME)-$$os-$$arch && cd ..; \
		fi; \
	done

# 运行测试
.PHONY: test
test:
	@echo "运行测试..."
	go test -v ./...

# 代码检查
.PHONY: lint
lint:
	@echo "运行代码检查..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint 未安装，跳过代码检查"; \
	fi

# 清理构建文件
.PHONY: clean
clean:
	@echo "清理构建文件..."
	@rm -rf bin/ dist/

# 安装到本地
.PHONY: install
install: build
	@echo "安装到本地..."
	cp bin/$(APP_NAME) $(GOPATH)/bin/$(APP_NAME)

# 创建发布
.PHONY: release
release: build-all
	@echo "准备发布 v$(VERSION)..."
	@echo "构建文件位于 dist/ 目录"

# 显示版本信息
.PHONY: version
version:
	@echo "$(APP_NAME) version $(VERSION)"
	@echo "Build time: $(BUILD_TIME)"
	@echo "Commit: $(COMMIT_HASH)"

# 帮助信息
.PHONY: help
help:
	@echo "可用命令："
	@echo "  build      - 构建当前平台"
	@echo "  build-all  - 构建所有平台"
	@echo "  test       - 运行测试"
	@echo "  lint       - 代码检查"
	@echo "  clean      - 清理构建文件"
	@echo "  install    - 安装到本地"
	@echo "  release    - 创建发布包"
	@echo "  version    - 显示版本信息"
	@echo "  help       - 显示此帮助" 