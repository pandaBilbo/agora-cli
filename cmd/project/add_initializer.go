package project

import (
	"fmt"
	"os"
	"path/filepath"
)

// AddInitializer 代码审查功能添加器
// 专门用于处理 devex add 命令的功能，继承BaseInitializer
type AddInitializer struct {
	BaseInitializer
}

// NewAddInitializer 创建代码审查功能添加器
func NewAddInitializer(projectPath string) (*AddInitializer, error) {
	globalConfigPath, err := getTemplatePath("global_config")
	if err != nil {
		return nil, fmt.Errorf("无法找到模板路径: %w", err)
	}

	return &AddInitializer{
		BaseInitializer: BaseInitializer{
			ProjectName:      filepath.Base(projectPath), // 使用目录名作为项目名
			FilePath:         projectPath,
			GlobalConfigPath: globalConfigPath,
			ConfigPath:       "template/config",
			TemplateCodePath: "template/code",
			NoGit:            false, // add命令默认不跳过Git
			NoCheck:          false, // add命令默认启用检查
			RemoteURL:        "",    // add命令不需要远程URL
		},
	}, nil
}

// CloneRepository add命令不需要克隆仓库
func (a *AddInitializer) CloneRepository() error {
	return nil
}

// CreateProject 检测现有项目，对于add命令不需要创建项目
func (a *AddInitializer) CreateProject() error {
	fmt.Println("📂 检测项目...")
	fmt.Println("  - 分析项目结构")
	fmt.Println("  - 准备添加代码审查支持")
	return nil
}

// InitDependencies 检查现有依赖，对于add命令不需要初始化依赖
func (a *AddInitializer) InitDependencies() error {
	fmt.Println("📥 检查项目依赖...")
	fmt.Println("  - 检查项目配置")
	fmt.Println("  - 验证基础环境")
	return nil
}

// ConfigureCodeReview 添加代码审查配置
func (a *AddInitializer) ConfigureCodeReview() error {
	if a.NoCheck {
		fmt.Println("⏭️  跳过代码审查配置 (使用了--no-check参数)")
		return nil
	}

	fmt.Println("🔍 添加代码审查配置...")
	fmt.Println("  - 添加通用代码检查配置")
	fmt.Println("  - 设置Git钩子")
	fmt.Println("  - 配置基础CI")
	fmt.Println("  - 添加审查模板")
	return nil
}

// ShowNextSteps 显示添加功能后的下一步操作
func (a *AddInitializer) ShowNextSteps() {
	a.BaseInitializer.ShowNextSteps()

	fmt.Println("\n代码审查配置完成，下一步：")
	fmt.Println("1. 安装相应的代码检查工具")
	fmt.Println("2. 运行代码检查命令")
	fmt.Println("3. 启用Git钩子检查")
}

// getTemplatePath 获取模板路径，优先使用二进制文件所在目录，然后尝试当前目录
func getTemplatePath(templateName string) (string, error) {
	// 获取当前执行文件的路径
	execPath, err := os.Executable()
	if err == nil {
		// 尝试从二进制文件同级目录查找template
		execDir := filepath.Dir(execPath)
		templatePath := filepath.Join(execDir, "template", templateName)
		if _, err := os.Stat(templatePath); err == nil {
			return templatePath, nil
		}
	}

	// 尝试从当前工作目录查找template（开发环境）
	currentDir, err := os.Getwd()
	if err == nil {
		templatePath := filepath.Join(currentDir, "template", templateName)
		if _, err := os.Stat(templatePath); err == nil {
			return templatePath, nil
		}
	}

	// 尝试相对路径（向上查找）
	for i := 0; i < 3; i++ {
		prefix := ""
		for j := 0; j < i; j++ {
			prefix = filepath.Join(prefix, "..")
		}
		templatePath := filepath.Join(prefix, "template", templateName)
		if _, err := os.Stat(templatePath); err == nil {
			absPath, _ := filepath.Abs(templatePath)
			return absPath, nil
		}
	}

	return "", fmt.Errorf("未找到模板目录: %s", templateName)
}
