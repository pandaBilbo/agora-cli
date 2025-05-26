# DevEx CLI

一个强大的开发者体验命令行工具。

## 安装

### 一键安装脚本（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/pandaBilbo/agora-cli/main/install.sh | bash
```

### 手动下载

从 [Releases 页面](https://github.com/pandaBilbo/agora-cli/releases) 下载适合你系统的二进制文件。

## 验证安装

```bash
devex version
```

## 使用

### 为现有项目添加代码审查功能

```bash
# 在项目根目录执行
devex add
```

这会为你的项目添加：
- 代码风格检查配置
- 敏感信息检查工具（gitleaks）
- Git钩子自动检查
- 代码审查模板

### 通过远程仓库初始化项目

```bash
devex init --remote https://github.com/username/your-repo.git
```

## 功能特性

- ✅ **全局安装支持** - 支持通过安装脚本全局安装
- ✅ **模板系统** - 内置多种项目模板
- ✅ **代码审查** - 自动配置代码检查工具
- ✅ **Git钩子** - 自动安装pre-commit钩子
- ✅ **多平台支持** - 支持macOS、Linux、Windows

## 开发

### 本地构建

```bash
make build
```

### 构建所有平台

```bash
make build-all
```

### 发布新版本

```bash
# 1. 构建所有平台
make build-all

# 2. 创建标签并推送
git tag v1.x.x
git push origin v1.x.x

# 3. 创建GitHub Release
gh release create v1.x.x dist/*.tar.gz --title "DevEx CLI v1.x.x" --notes "发布说明"
```