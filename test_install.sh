#!/bin/bash

# DevEx CLI 安装脚本
set -e

# 默认值
REPO="pandaBilbo/agora-cli"  # 请替换为您的实际仓库
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="devex"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

print_green() {
    echo -e "${GREEN}$1${NC}"
}

print_yellow() {
    echo -e "${YELLOW}$1${NC}"
}

print_red() {
    echo -e "${RED}$1${NC}"
}

# 获取系统信息
get_os() {
    case "$(uname -s)" in
        Darwin) echo "darwin" ;;
        Linux) echo "linux" ;;
        CYGWIN*|MINGW32*|MSYS*|MINGW*) echo "windows" ;;
        *) 
            print_red "不支持的操作系统: $(uname -s)"
            exit 1
            ;;
    esac
}

get_arch() {
    case "$(uname -m)" in
        x86_64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *)
            print_red "不支持的架构: $(uname -m)"
            exit 1
            ;;
    esac
}

# 获取最新版本
get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# 主安装函数
install_devex() {
    print_green "🚀 开始安装 DevEx CLI..."
    
    OS=$(get_os)
    ARCH=$(get_arch)
    echo "DEBUG: Trying to get version..."; VERSION=$(get_latest_version); echo "DEBUG: Got version: $VERSION"
    
    if [ -z "$VERSION" ]; then
        print_red "无法获取最新版本信息"
        exit 1
    fi
    
    print_yellow "检测到系统: $OS/$ARCH"
    print_yellow "最新版本: $VERSION"
    
    # 构建下载URL
    FILENAME="${BINARY_NAME}-${OS}-${ARCH}"
    if [ "$OS" = "windows" ]; then
        FILENAME="${FILENAME}.exe"
    fi
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}-${OS}-${ARCH}.tar.gz"
    
    print_yellow "下载地址: $DOWNLOAD_URL"
    
    # 创建临时目录
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    # 下载并解压
    print_yellow "正在下载 $BINARY_NAME $VERSION..."
    if ! curl -L "$DOWNLOAD_URL" -o "${BINARY_NAME}.tar.gz"; then
        print_red "下载失败"
        exit 1
    fi
    
    print_yellow "正在解压..."
    tar -xzf "${BINARY_NAME}.tar.gz"
    
    # 安装
    print_yellow "正在安装到 $INSTALL_DIR..."
    
    BINARY_PATH="${BINARY_NAME}-${OS}-${ARCH}/${BINARY_NAME}"
    if [ "$OS" = "windows" ]; then
        BINARY_PATH="${BINARY_NAME}-${OS}-${ARCH}/${BINARY_NAME}.exe"
    fi
    
    if [ ! -f "$BINARY_PATH" ]; then
        print_red "二进制文件不存在: $BINARY_PATH"
        exit 1
    fi
    
    # 检查权限并安装
    if [ -w "$INSTALL_DIR" ]; then
        cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    else
        print_yellow "需要管理员权限安装到 $INSTALL_DIR"
        sudo cp "$BINARY_PATH" "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # 设置可执行权限
    if [ "$OS" != "windows" ]; then
        if [ -w "$INSTALL_DIR/$BINARY_NAME" ]; then
            chmod +x "$INSTALL_DIR/$BINARY_NAME"
        else
            sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
        fi
    fi
    
    # 清理
    cd - > /dev/null
    rm -rf "$TMP_DIR"
    
    print_green "✅ DevEx CLI $VERSION 安装成功!"
    print_green "验证安装: $BINARY_NAME version"
    
    # 验证安装
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        print_green "$("$BINARY_NAME" version)"
    else
        print_yellow "请将 $INSTALL_DIR 添加到您的 PATH 环境变量中"
    fi
}

# 检查依赖
check_dependencies() {
    if ! command -v curl >/dev/null 2>&1; then
        print_red "需要 curl 命令"
        exit 1
    fi
    
    if ! command -v tar >/dev/null 2>&1; then
        print_red "需要 tar 命令"
        exit 1
    fi
}

# 显示帮助
show_help() {
    cat << EOF
DevEx CLI 安装脚本

用法:
    $0 [选项]

选项:
    -h, --help          显示此帮助信息
    -d, --dir DIR       指定安装目录 (默认: $INSTALL_DIR)
    --repo REPO         指定GitHub仓库 (默认: $REPO)

示例:
    # 默认安装
    $0
    
    # 安装到指定目录
    $0 --dir ~/bin
    
    # 从指定仓库安装
    $0 --repo username/repo-name

EOF
}

# 解析命令行参数
while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -d|--dir)
            INSTALL_DIR="$2"
            shift 2
            ;;
        --repo)
            REPO="$2"
            shift 2
            ;;
        *)
            print_red "未知选项: $1"
            show_help
            exit 1
            ;;
    esac
done

# 主流程
main() {
    check_dependencies
    install_devex
}

main "$@" 