package project

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Initializer 定义项目初始化器的接口
type Initializer interface {
	// CloneRepository 克隆远程仓库
	CloneRepository() error

	// CreateProject 创建项目特定文件
	CreateProject() error

	// CopyTemplateFiles 复制模板文件
	CopyTemplateFiles() error

	// InitDependencies 初始化依赖
	InitDependencies() error

	// InstallGitHooks 安装 Git 钩子
	InstallGitHooks() error

	// ConfigureCodeReview 配置代码审查
	ConfigureCodeReview() error

	// ShowNextSteps 显示后续步骤
	ShowNextSteps()
}

// BaseInitializer 提供基础实现
type BaseInitializer struct {
	ProjectName      string
	FilePath         string
	GlobalConfigPath string
	ConfigPath       string
	TemplateCodePath string
	NoGit            bool
	NoCheck          bool
	RemoteURL        string
}

// CloneRepository 克隆远程仓库的基础实现
func (b *BaseInitializer) CloneRepository() error {
	fmt.Printf("📥 克隆远程仓库: %s\n", b.RemoteURL)
	if b.RemoteURL == "" {
		return fmt.Errorf("未指定远程仓库地址")
	}

	// 检查git命令是否存在
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("未找到git命令，请先安装git")
	}

	// 检查目标目录是否已存在
	if _, err := os.Stat(b.FilePath); !os.IsNotExist(err) {
		return fmt.Errorf("目标目录已存在: %s", b.FilePath)
	}

	// 执行git clone
	fmt.Printf("  - 正在克隆到: %s\n", b.FilePath)
	cmd := exec.Command("git", "clone", b.RemoteURL, b.FilePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("克隆远程仓库失败: %w", err)
	}

	fmt.Println("  ✅ 仓库克隆成功")
	return nil
}

// CreateProject 创建项目特定文件的基础实现
func (b *BaseInitializer) CreateProject() error {
	// 基础实现不做任何事情
	return nil
}

// CopyTemplateFiles 复制模板文件的基础实现
func (b *BaseInitializer) CopyTemplateFiles() error {
	fmt.Println("📂 复制模板文件...")

	if err := copyDir(b.GlobalConfigPath, b.FilePath); err != nil {
		return fmt.Errorf("复制全局配置文件失败: %w", err)
	}
	fmt.Printf("  - 已复制全局配置文件: %s\n", b.GlobalConfigPath)

	fmt.Println("  ✅ 模板文件复制完成")
	return nil
}

// InitDependencies 初始化依赖的基础实现
func (b *BaseInitializer) InitDependencies() error {
	return nil
}

// InstallGitHooks 安装 Git 钩子的基础实现
func (b *BaseInitializer) InstallGitHooks() error {
	if b.NoCheck {
		fmt.Println("⏭️  跳过Git钩子安装 (使用了--no-check参数)")
		return nil
	}

	fmt.Println("🔗 安装Git钩子...")

	hookScript := filepath.Join(b.FilePath, ".git-hooks", "install-hooks.sh")

	// 检查脚本是否存在
	if _, err := os.Stat(hookScript); os.IsNotExist(err) {
		fmt.Println("  ⚠️  Git钩子脚本不存在，跳过安装")
		return nil
	}

	// 获取绝对路径
	absHookScript, err := filepath.Abs(hookScript)
	if err != nil {
		return fmt.Errorf("获取脚本绝对路径失败: %w", err)
	}

	cmd := exec.Command("bash", absHookScript)
	cmd.Dir = b.FilePath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("安装Git钩子失败: %w", err)
	}

	fmt.Println("  ✅ Git钩子安装成功")
	return nil
}

// ConfigureCodeReview 配置代码审查的基础实现
func (b *BaseInitializer) ConfigureCodeReview() error {
	if b.NoCheck {
		return nil
	}
	return nil
}

// ShowNextSteps 显示后续步骤的基础实现
func (b *BaseInitializer) ShowNextSteps() {
	fmt.Println("\n✨ 项目创建成功！")
	fmt.Printf("进入项目目录：cd %s\n", b.ProjectName)
}

// copyDir 递归复制目录
func copyDir(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// copyFile 复制单个文件
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	// 保持文件权限
	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return err
	}

	return nil
}

// getTemplatePath 获取模板路径，优先使用二进制文件所在目录，然后尝试当前目录
func getTemplatePath(templateName string) (string, error) {
	// 获取当前执行文件的路径
	execPath, err := os.Executable()
	if err == nil {
		execDir := filepath.Dir(execPath)

		// 尝试从二进制文件同级目录查找template
		templatePath := filepath.Join(execDir, "template", templateName)
		if _, err := os.Stat(templatePath); err == nil {
			return templatePath, nil
		}

		// 尝试从二进制文件上级目录查找template（开发环境）
		templatePath = filepath.Join(execDir, "..", "template", templateName)
		if _, err := os.Stat(templatePath); err == nil {
			absPath, _ := filepath.Abs(templatePath)
			return absPath, nil
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
