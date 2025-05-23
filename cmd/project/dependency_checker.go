package project

import (
	"fmt"
	"os/exec"
	"strings"
)

// DependencyChecker 依赖检查器接口
type DependencyChecker interface {
	CheckDependencies(dependencies []string) error
	CheckSingleDependency(command string) error
	GetMissingDependencies(dependencies []string) []string
}

// CommandDependencyChecker 命令行工具依赖检查器
type CommandDependencyChecker struct{}

// NewCommandDependencyChecker 创建命令行依赖检查器
func NewCommandDependencyChecker() *CommandDependencyChecker {
	return &CommandDependencyChecker{}
}

// CheckDependencies 检查多个依赖
func (c *CommandDependencyChecker) CheckDependencies(dependencies []string) error {
	missing := c.GetMissingDependencies(dependencies)
	if len(missing) > 0 {
		return fmt.Errorf("缺少必需的命令行工具: %s", strings.Join(missing, ", "))
	}
	return nil
}

// CheckSingleDependency 检查单个依赖
func (c *CommandDependencyChecker) CheckSingleDependency(command string) error {
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("未找到命令: %s", command)
	}
	return nil
}

// GetMissingDependencies 获取缺失的依赖列表
func (c *CommandDependencyChecker) GetMissingDependencies(dependencies []string) []string {
	var missing []string
	for _, dep := range dependencies {
		if err := c.CheckSingleDependency(dep); err != nil {
			missing = append(missing, dep)
		}
	}
	return missing
}

// SwiftDependencyHelper Swift特定的依赖检查帮助器
type SwiftDependencyHelper struct {
	checker DependencyChecker
}

// NewSwiftDependencyHelper 创建Swift依赖检查帮助器
func NewSwiftDependencyHelper() *SwiftDependencyHelper {
	return &SwiftDependencyHelper{
		checker: NewCommandDependencyChecker(),
	}
}

// CheckAndInstallXcodegen 检查并安装xcodegen
func (h *SwiftDependencyHelper) CheckAndInstallXcodegen() error {
	// 检查是否已安装
	if err := h.checker.CheckSingleDependency("xcodegen"); err == nil {
		return nil
	}

	// 未安装，检查是否有 Homebrew
	if err := h.checker.CheckSingleDependency("brew"); err != nil {
		return fmt.Errorf("未找到 Homebrew，请先安装 Homebrew：\n/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
	}

	fmt.Println("正在安装 xcodegen...")

	// 执行安装命令
	cmd := exec.Command("brew", "install", "xcodegen")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("安装 xcodegen 失败：%w", err)
	}

	fmt.Println("xcodegen 安装成功！")
	return nil
}

// GetInstallationInstructions 获取依赖安装说明
func GetInstallationInstructions(command string) string {
	instructions := map[string]string{
		"xcodegen": "安装说明：brew install xcodegen",
		"pod":      "安装说明：sudo gem install cocoapods",
		"brew":     "安装说明：/bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"",
		"go":       "安装说明：https://golang.org/dl/",
		"java":     "安装说明：https://adoptopenjdk.net/",
		"gradle":   "安装说明：https://gradle.org/install/",
		"mvn":      "安装说明：https://maven.apache.org/install.html",
	}

	if instruction, exists := instructions[command]; exists {
		return instruction
	}
	return fmt.Sprintf("请手动安装 %s", command)
}
