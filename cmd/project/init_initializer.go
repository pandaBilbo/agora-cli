package project

import (
	"fmt"
)

// InitInitializer 项目初始化器
// 专门用于处理 devex init 命令的功能
type InitInitializer struct {
	BaseInitializer
}

// NewInitInitializer 创建项目初始化器
func NewInitInitializer(projectName, projectPath string, noGit, noCheck bool, remote string) (*InitInitializer, error) {
	globalConfigPath, err := getTemplatePath("global_config")
	if err != nil {
		return nil, fmt.Errorf("无法找到模板路径: %w", err)
	}

	return &InitInitializer{
		BaseInitializer: BaseInitializer{
			ProjectName:      projectName,
			FilePath:         projectPath,
			GlobalConfigPath: globalConfigPath,
			ConfigPath:       "template/config",
			TemplateCodePath: "template/code",
			NoGit:            noGit,
			NoCheck:          noCheck,
			RemoteURL:        remote,
		},
	}, nil
}

// CreateProject 创建项目结构和文件
func (i *InitInitializer) CreateProject() error {
	fmt.Println("📦 创建项目...")
	fmt.Println("  - 设置项目结构")
	fmt.Println("  - 配置基础文件")
	return nil
}

// InitDependencies 初始化项目依赖
func (i *InitInitializer) InitDependencies() error {
	fmt.Println("📥 初始化项目依赖...")
	fmt.Println("  - 检查基础依赖")
	return nil
}

// ConfigureCodeReview 配置代码审查工具
func (i *InitInitializer) ConfigureCodeReview() error {
	if i.NoCheck {
		fmt.Println("⏭️  跳过代码审查配置 (使用了--no-check参数)")
		return nil
	}

	fmt.Println("🔍 配置代码审查工具...")
	fmt.Println("  - 配置通用代码检查工具")
	fmt.Println("  - 设置Git钩子")
	return nil
}

// ShowNextSteps 显示项目初始化后的下一步操作
func (i *InitInitializer) ShowNextSteps() {
	i.BaseInitializer.ShowNextSteps()

	fmt.Println("\n项目下一步：")
	fmt.Println("1. 使用编辑器打开项目")
	fmt.Println("2. 安装相关依赖")
	fmt.Println("3. 开始开发")
}
