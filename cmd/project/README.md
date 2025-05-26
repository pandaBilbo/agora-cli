# 项目初始化器架构

这个包提供了一个可扩展的项目初始化器架构，支持多种编程语言的项目创建。

## 🏗️ 架构概览

```
cmd/project/
├── initializer.go      # 基础接口和通用实现
├── template_manager.go # 模板管理系统
├── language_config.go  # 语言配置管理
├── factory.go          # 初始化器工厂
├── swift.go           # Swift 特定实现
└── README.md          # 本文档
```

## 🚀 如何添加新语言支持

添加新语言支持非常简单，只需要3个步骤：

### 1. 在 `language_config.go` 中添加语言配置

```go
// 在 getLanguageConfigs() 函数中添加新语言
"go": {
    Name:             "go",
    DisplayName:      "Go",
    TemplateCodePath: filepath.Join("template", "go", "code"),
    ConfigPath:       filepath.Join("template", "go", "config"),
    GlobalConfigPath: filepath.Join("template", "global_config"),
    RequiredCommands: []string{"go"},
    FileExtensions:   []string{".go"},
    BuildTool:        "go",
    PackageManager:   "Go modules",
},
```

### 2. 在 `factory.go` 中添加初始化器创建逻辑

```go
// 在 NewInitializer() 函数的 switch 语句中添加
case "go":
    return NewGoInitializer(projectName, path, config.GlobalConfigPath, config.ConfigPath, config.TemplateCodePath, noGit, noCheck, remote), nil
```

### 3. 创建语言特定的初始化器

创建 `go_initializer.go` 文件，实现 `Initializer` 接口：

```go
type GoInitializer struct {
    BaseInitializer
    templates TemplateManager
}

func NewGoInitializer(...) *GoInitializer { ... }
func (g *GoInitializer) CreateProject() error { ... }
// 实现其他接口方法...
```

### 4. 准备模板文件

创建对应的模板目录结构：
```
template/
├── global_config/     # 全局配置文件
├── go/
│   ├── config/       # Go 特定配置
│   └── code/         # Go 代码模板
│       ├── main.go
│       ├── go.mod.tpl
│       └── ...
```

## 🛠️ 开发环境

### 环境要求

- Go 1.21+
- Make
- Git
- GitHub CLI (gh) - 用于发布

### 本地开发

```bash
# 克隆项目
git clone https://github.com/pandaBilbo/agora-cli.git
cd agora-cli

# 本地构建
make build

# 运行测试
make test

# 代码检查
make lint

# 清理构建文件
make clean
```

### 调试

```bash
# 使用本地构建的版本测试
./bin/devex init --remote https://github.com/username/test-repo.git

# 调试add命令
./bin/devex add --path /path/to/project
```

## 🚀 发布流程

### 1. 准备发布

```bash
# 确保所有更改已提交
git add .
git commit -m "feat: 添加新功能"
git push origin main

# 确保本地是最新的
git pull origin main
```

### 2. 版本发布

```bash
# 创建并推送标签
git tag v1.x.x
git push origin v1.x.x

# 构建所有平台的二进制文件
make release

# 创建GitHub Release
gh release create v1.x.x \
  --title "DevEx CLI v1.x.x" \
  --notes "发布说明：
  - 新增功能：xxx
  - 修复问题：xxx
  - 性能优化：xxx"

# 上传构建产物
gh release upload v1.x.x dist/*.tar.gz
```

### 3. 发布后验证

```bash
# 测试安装脚本
curl -fsSL https://raw.githubusercontent.com/pandaBilbo/agora-cli/main/install.sh | bash

# 验证版本
devex version

# 测试核心功能
devex init --remote https://github.com/pandaBilbo/cli.git
```

## 🔧 Make 命令

```bash
# 构建当前平台
make build

# 构建所有平台
make build-all

# 创建发布包（包含template目录）
make release

# 运行测试
make test

# 代码检查
make lint

# 清理构建文件
make clean

# 安装到本地GOPATH
make install

# 显示版本信息
make version

# 显示帮助
make help
```

## 📦 发布目录结构

构建完成后，`dist/` 目录结构如下：

```
dist/
├── devex-darwin-amd64/
│   ├── devex
│   └── template/
├── devex-darwin-amd64.tar.gz
├── devex-darwin-arm64/
│   ├── devex
│   └── template/
├── devex-darwin-arm64.tar.gz
├── devex-linux-amd64/
│   ├── devex
│   └── template/
├── devex-linux-amd64.tar.gz
├── devex-linux-arm64/
│   ├── devex
│   └── template/
├── devex-linux-arm64.tar.gz
├── devex-windows-amd64/
│   ├── devex.exe
│   └── template/
└── devex-windows-amd64.tar.gz
```

## 🎯 设计原则

1. **简单性**：避免过度设计，保持代码简洁
2. **可扩展性**：新语言支持只需要最少的代码修改
3. **一致性**：所有语言使用相同的接口和模式
4. **用户友好**：提供清晰的错误信息和进度提示

## 🔧 模板系统

模板管理器支持：
- 模板文件存在性检查
- 变量替换（支持 `${var}` 和 `{{var}}` 格式）
- 模板列表查询
- 简单的错误处理

## 📝 最佳实践

1. **错误处理**：使用 `fmt.Errorf` 包装错误，提供上下文信息
2. **进度提示**：在长时间操作时显示进度信息
3. **模板检查**：在使用模板前检查文件是否存在
4. **配置验证**：验证必需的命令行工具是否可用
5. **版本管理**：使用语义化版本控制
6. **持续集成**：每次发布前运行完整测试

## 🧪 测试

```bash
# 编译测试
go build ./cmd/project

# 运行测试（如果有的话）
go test ./cmd/project

# 集成测试
./test/integration_test.sh
```

## 🔍 故障排除

### 构建问题

```bash
# 清理并重新构建
make clean && make build

# 检查Go版本
go version

# 检查依赖
go mod verify
```

### 发布问题

```bash
# 检查GitHub CLI配置
gh auth status

# 检查标签是否存在
git tag -l

# 重新发布（删除已存在的release）
gh release delete v1.x.x
git tag -d v1.x.x
git push origin :refs/tags/v1.x.x
```

## 🎯 设计原则

1. **简单性**：避免过度设计，保持代码简洁
2. **可扩展性**：新语言支持只需要最少的代码修改
3. **一致性**：所有语言使用相同的接口和模式
4. **用户友好**：提供清晰的错误信息和进度提示

## 🔧 模板系统

模板管理器支持：
- 模板文件存在性检查
- 变量替换（支持 `${var}` 和 `{{var}}` 格式）
- 模板列表查询
- 简单的错误处理

## 📝 最佳实践

1. **错误处理**：使用 `fmt.Errorf` 包装错误，提供上下文信息
2. **进度提示**：在长时间操作时显示进度信息
3. **模板检查**：在使用模板前检查文件是否存在
4. **配置验证**：验证必需的命令行工具是否可用

## 🧪 测试

```bash
# 编译测试
go build ./cmd/project

# 运行测试（如果有的话）
go test ./cmd/project
``` 