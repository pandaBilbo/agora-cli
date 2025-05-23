# DevEx CLI 发布指南

## 快速发布

### 1. 手动发布（推荐用于测试）

```bash
# 构建当前平台
make build

# 构建所有平台
make build-all

# 查看构建文件
ls -la dist/
```

### 2. 自动发布（推荐用于正式发布）

创建并推送Git标签即可触发自动发布：

```bash
# 创建标签
git tag v1.0.0

# 推送标签到远程仓库
git push origin v1.0.0
```

或者通过GitHub Actions手动触发：
1. 访问GitHub仓库的Actions页面
2. 选择"Release"工作流
3. 点击"Run workflow"
4. 输入版本号（如 v1.0.0）
5. 点击"Run workflow"

## 发布流程详解

### 版本号规范

采用[语义化版本](https://semver.org/lang/zh-CN/)规范：

- **主版本号**：当你做了不兼容的API修改
- **次版本号**：当你做了向下兼容的功能性新增
- **修订号**：当你做了向下兼容的问题修正

示例：
- `v1.0.0` - 第一个稳定版本
- `v1.1.0` - 添加新功能
- `v1.1.1` - 修复Bug
- `v2.0.0` - 重大更新，可能不向下兼容

### 发布前检查清单

- [ ] 代码已提交并推送到主分支
- [ ] 所有测试通过：`make test`
- [ ] 代码检查通过：`make lint`
- [ ] 版本号已更新
- [ ] 更新日志已准备

### 本地构建测试

```bash
# 清理之前的构建
make clean

# 构建并测试
make build
./bin/devex version

# 构建所有平台
make build-all

# 检查构建结果
ls -la dist/
```

### 创建发布

#### 方法1：通过Git标签（推荐）

```bash
# 确保在主分支且代码是最新的
git checkout main
git pull origin main

# 创建并推送标签
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

#### 方法2：通过GitHub界面

1. 访问GitHub仓库
2. 点击"Releases"
3. 点击"Create a new release"
4. 填写标签和版本信息
5. 发布

#### 方法3：通过GitHub Actions手动触发

1. 访问Actions页面
2. 选择"Release"工作流
3. 点击"Run workflow"
4. 输入版本号
5. 运行

## 安装方式

### 自动安装（推荐）

用户可以使用一键安装脚本：

```bash
curl -fsSL https://raw.githubusercontent.com/你的用户名/agora-cli/main/install.sh | bash
```

### 手动安装

#### macOS

```bash
# Intel芯片
curl -L https://github.com/你的用户名/agora-cli/releases/latest/download/devex-darwin-amd64.tar.gz | tar xz
sudo mv devex-darwin-amd64/devex /usr/local/bin/

# Apple Silicon
curl -L https://github.com/你的用户名/agora-cli/releases/latest/download/devex-darwin-arm64.tar.gz | tar xz
sudo mv devex-darwin-arm64/devex /usr/local/bin/
```

#### Linux

```bash
curl -L https://github.com/你的用户名/agora-cli/releases/latest/download/devex-linux-amd64.tar.gz | tar xz
sudo mv devex-linux-amd64/devex /usr/local/bin/
```

#### Windows

1. 访问[Releases页面](https://github.com/你的用户名/agora-cli/releases)
2. 下载`devex-windows-amd64.tar.gz`
3. 解压并将`devex.exe`添加到PATH

### Go用户

```bash
go install github.com/你的用户名/agora-cli@latest
```

## 发布后任务

### 1. 验证发布

```bash
# 测试安装脚本
curl -fsSL https://raw.githubusercontent.com/你的用户名/agora-cli/main/install.sh | bash

# 验证版本
devex version
```

### 2. 更新文档

- [ ] 更新README.md中的安装说明
- [ ] 更新CHANGELOG.md
- [ ] 更新相关博客或文档

### 3. 社区通知

- [ ] 在相关社区发布更新通知
- [ ] 更新包管理器（如有）

## 常见问题

### 发布失败

1. **权限问题**：确保GitHub Token有正确的权限
2. **构建失败**：检查Go版本和依赖
3. **标签冲突**：删除本地和远程的冲突标签后重新创建

### 撤回发布

```bash
# 删除远程标签
git push --delete origin v1.0.0

# 删除本地标签
git tag --delete v1.0.0

# 在GitHub上删除Release
```

### 修复发布

如果发布有问题但不想删除：

1. 创建新的补丁版本（如v1.0.1）
2. 在GitHub Releases中标记有问题的版本为"Pre-release"
3. 在Release描述中说明问题

## 自动化改进

### 添加Homebrew支持

创建Homebrew Formula，让用户可以通过`brew install devex`安装。

### 添加更多包管理器

- **Scoop** (Windows)
- **Chocolatey** (Windows)  
- **AUR** (Arch Linux)
- **Snap** (Ubuntu)

### CI/CD增强

- 添加自动测试
- 添加安全扫描
- 添加性能基准测试
- 添加多平台测试 