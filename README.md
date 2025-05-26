# DevEx CLI

一个强大的开发者体验命令行工具，帮助你快速为项目添加代码审查和质量检查功能。

## 快速安装

### 一键安装（推荐）

```bash
curl -fsSL https://raw.githubusercontent.com/pandaBilbo/agora-cli/main/install.sh | bash
```

### 手动下载

从 [Releases 页面](https://github.com/pandaBilbo/agora-cli/releases) 下载适合你系统的二进制文件。

### 验证安装

```bash
devex version
```

## 使用方法

### 为现有项目添加代码质量检查

```bash
# 在项目根目录执行
devex add
```

这会为你的项目添加：
- 代码风格检查配置
- 敏感信息泄露检测
- Git提交钩子
- 代码审查模板

### 通过远程仓库初始化项目

```bash
devex init --remote https://github.com/username/your-repo.git
```

### 查看帮助

```bash
devex --help
devex add --help
devex init --help
```

## 功能特性

- ✅ **一键安装** - 支持macOS、Linux、Windows
- ✅ **自动依赖管理** - 自动检测并安装pre-commit和gitleaks
- ✅ **代码质量检查** - pre-commit钩子自动检查代码风格
- ✅ **敏感信息保护** - 集成gitleaks防止密钥泄露
- ✅ **提交信息规范** - 防止提交信息包含中文字符
- ✅ **模板系统** - 快速初始化项目配置

## 故障排除

### 安装失败

如果安装脚本失败，请检查：

1. 网络连接是否正常
2. 是否有sudo权限（需要写入/usr/local/bin）

也可以手动下载并安装：

```bash
# Linux示例
wget https://github.com/pandaBilbo/agora-cli/releases/latest/download/devex-linux-amd64.tar.gz
tar -xzf devex-linux-amd64.tar.gz
sudo cp devex-linux-amd64/devex /usr/local/bin/
sudo cp -r devex-linux-amd64/template /usr/local/bin/
```

### Git钩子安装失败

DevEx CLI会自动检测并安装所需依赖。如果自动安装失败，可以手动安装：

```bash
# 手动安装pre-commit（通常自动安装）
pip install pre-commit

# 手动安装gitleaks（macOS会自动通过Homebrew安装）
brew install gitleaks
```

**注意**：大多数情况下不需要手动安装，工具会自动处理依赖。

## 支持

- 📝 [提交Issue](https://github.com/pandaBilbo/agora-cli/issues)
- 💬 [讨论区](https://github.com/pandaBilbo/agora-cli/discussions)

## 许可证

MIT License