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