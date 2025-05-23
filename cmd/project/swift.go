package project

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SwiftInitializer Swift项目初始化器
// 使用 TemplateManager 替代硬编码模板
type SwiftInitializer struct {
	BaseInitializer
	templates        TemplateManager
	dependencyHelper *SwiftDependencyHelper
	config           *LanguageConfig
}

// NewSwiftInitializer 创建Swift项目初始化器
func NewSwiftInitializer(projectName, projectPath, globalConfigPath, configPath, templateCodePath string, noGit, noCheck bool, remote string) *SwiftInitializer {
	templateManager := NewFileTemplateManager(templateCodePath)
	dependencyHelper := NewSwiftDependencyHelper()

	// 获取Swift配置
	config, _ := GetLanguageConfig("swift")

	return &SwiftInitializer{
		BaseInitializer: BaseInitializer{
			ProjectName:      projectName,
			FilePath:         projectPath,
			GlobalConfigPath: globalConfigPath,
			ConfigPath:       configPath,
			TemplateCodePath: templateCodePath,
			NoGit:            noGit,
			NoCheck:          noCheck,
			RemoteURL:        remote,
		},
		templates:        templateManager,
		dependencyHelper: dependencyHelper,
		config:           config,
	}
}

// CopyTemplateFiles Swift项目特定的模板文件复制
func (s *SwiftInitializer) CopyTemplateFiles() error {
	// 先调用父类的方法复制所有文件
	if err := s.BaseInitializer.CopyTemplateFiles(); err != nil {
		return err
	}

	// 读取Podfile文件
	podfilePath := filepath.Join(s.FilePath, "Podfile")
	content, err := os.ReadFile(podfilePath)
	if err != nil {
		return fmt.Errorf("读取Podfile失败：%w", err)
	}

	// 替换target名称
	oldContent := string(content)
	newContent := strings.ReplaceAll(
		oldContent,
		"target 'XXXX' do",
		fmt.Sprintf("target '%s' do", s.ProjectName),
	)

	// 写回文件
	if err := os.WriteFile(podfilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("更新Podfile失败：%w", err)
	}

	// 创建项目文件
	if err := s.createProjectFiles(); err != nil {
		return err
	}

	return nil
}

// CreateProject Swift项目特定的项目创建
func (s *SwiftInitializer) CreateProject() error {
	// 检查依赖
	if err := s.checkDependencies(); err != nil {
		return err
	}

	fmt.Println("正在生成 Xcode 项目...")

	// 使用 xcodegen 生成项目文件
	cmd := exec.Command("xcodegen", "generate", "--spec", "project.yml")
	cmd.Dir = s.FilePath

	// 捕获命令的输出
	var stdout, stderr bytes.Buffer
	cmd.Stdout = io.MultiWriter(os.Stdout, &stdout)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stderr)

	if err := cmd.Run(); err != nil {
		// 如果命令执行失败，返回详细的错误信息
		return fmt.Errorf("生成 Xcode 项目失败：\n错误：%v\n输出：%s\n错误输出：%s",
			err, stdout.String(), stderr.String())
	}

	fmt.Println("✅ Xcode 项目生成成功！")
	return nil
}

// checkDependencies 检查所有依赖
func (s *SwiftInitializer) checkDependencies() error {
	fmt.Println("🔍 检查依赖...")

	checker := NewCommandDependencyChecker()
	missing := checker.GetMissingDependencies(s.config.RequiredCommands)

	if len(missing) == 0 {
		fmt.Println("  ✅ 所有依赖都已安装")
		return nil
	}

	fmt.Printf("  ⚠️  缺少依赖: %s\n", strings.Join(missing, ", "))

	// 尝试自动安装 xcodegen
	for _, cmd := range missing {
		if cmd == "xcodegen" {
			fmt.Println("  🔧 尝试自动安装 xcodegen...")
			if err := s.dependencyHelper.CheckAndInstallXcodegen(); err != nil {
				return err
			}
		} else {
			// 对于其他依赖，提供安装说明
			fmt.Printf("  💡 %s\n", GetInstallationInstructions(cmd))
		}
	}

	// 再次检查是否还有缺失的依赖
	stillMissing := checker.GetMissingDependencies(s.config.RequiredCommands)
	if len(stillMissing) > 0 {
		return fmt.Errorf("仍然缺少依赖: %s", strings.Join(stillMissing, ", "))
	}

	return nil
}

// createProjectFiles 创建项目所需的所有文件
func (s *SwiftInitializer) createProjectFiles() error {
	// 创建项目源代码目录
	srcDir := filepath.Join(s.FilePath, s.ProjectName)
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		return fmt.Errorf("创建源代码目录失败：%w", err)
	}

	// 定义所有需要创建的文件
	files := map[string]string{
		"Info.plist":              "Info.plist",
		"AppDelegate.swift":       "AppDelegate.swift",
		"ViewController.swift":    "ViewController.swift",
		"Main.storyboard":         "Main.storyboard",
		"LaunchScreen.storyboard": "LaunchScreen.storyboard",
	}

	// 检查模板文件是否存在
	for _, tplName := range files {
		if !s.templates.TemplateExists(tplName) {
			return fmt.Errorf("模板文件不存在: %s。请检查模板目录: %s", tplName, s.TemplateCodePath)
		}
	}

	fmt.Println("📄 创建项目文件...")

	// 写入所有文件
	for filename, tplName := range files {
		fmt.Printf("  - 创建 %s\n", filename)
		content, err := s.templates.LoadTemplateCode(tplName)
		if err != nil {
			return fmt.Errorf("加载模板 %s 失败：%w", tplName, err)
		}
		filePath := filepath.Join(srcDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return fmt.Errorf("创建文件 %s 失败：%w", filename, err)
		}
	}

	// 检查 project.yml 模板
	if !s.templates.TemplateExists("project.yml") {
		return fmt.Errorf("模板文件不存在: project.yml。请检查模板目录: %s", s.TemplateCodePath)
	}

	// 写入 project.yml 文件，替换变量
	fmt.Println("  - 创建 project.yml")
	content, err := s.templates.RenderTemplateCode("project.yml", map[string]string{"PROJECT_NAME": s.ProjectName})
	if err != nil {
		return fmt.Errorf("渲染 project.yml 失败：%w", err)
	}
	filePath := filepath.Join(s.FilePath, "project.yml")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("创建 project.yml 失败：%w", err)
	}
	return nil
}

// InitDependencies 初始化依赖，如pod install
func (s *SwiftInitializer) InitDependencies() error {
	// 显示代码审查工具说明（仅在启用代码审查时）
	if !s.NoCheck {
		fmt.Println("💡 代码审查工具说明：")
		fmt.Println("  - SwiftLint: 代码风格检查，会在构建时自动运行")
		fmt.Println("  - gitleaks: 敏感信息检查，通过Git钩子运行")
		fmt.Println("  - 手动运行SwiftLint: swiftlint lint")
	}

	return nil
}

// ConfigureCodeReview Swift项目特定的代码审查配置
func (s *SwiftInitializer) ConfigureCodeReview() error {
	if s.NoCheck {
		fmt.Println("⏭️  跳过代码审查配置 (使用了--no-check参数)")
		return nil
	}

	fmt.Println("🔍 配置Swift代码审查工具...")

	// 1. 检查SwiftLint配置文件是否存在
	swiftlintPath := filepath.Join(s.FilePath, ".swiftlint.yml")
	if _, err := os.Stat(swiftlintPath); os.IsNotExist(err) {
		fmt.Printf("  ⚠️  SwiftLint配置文件不存在: %s\n", swiftlintPath)
		fmt.Println("  💡 提示：配置文件应该在复制模板文件时创建")
	} else {
		fmt.Println("  ✅ SwiftLint配置文件已存在")
	}

	// 2. 向Podfile添加SwiftLint依赖
	if err := s.addSwiftLintToPodfile(); err != nil {
		return fmt.Errorf("添加SwiftLint依赖失败: %w", err)
	}

	// 3. 向project.yml添加SwiftLint构建脚本
	if err := s.addSwiftLintToProjectYml(); err != nil {
		return fmt.Errorf("添加SwiftLint脚本失败: %w", err)
	}

	fmt.Println("  ✅ 代码审查配置完成")
	fmt.Println("  💡 提示：依赖将在下一步安装")

	return nil
}

// addSwiftLintToPodfile 向Podfile添加SwiftLint依赖
func (s *SwiftInitializer) addSwiftLintToPodfile() error {
	podfilePath := filepath.Join(s.FilePath, "Podfile")

	// 读取现有Podfile内容
	content, err := os.ReadFile(podfilePath)
	if err != nil {
		return fmt.Errorf("读取Podfile失败: %w", err)
	}

	podfileContent := string(content)

	// 检查是否已经包含SwiftLint依赖
	if strings.Contains(podfileContent, "SwiftLint") {
		fmt.Println("  ✅ Podfile已包含SwiftLint依赖")
		return nil
	}

	// 在target块内添加SwiftLint依赖
	// 查找"use_frameworks!"后面的位置
	useFrameworksPos := strings.Index(podfileContent, "use_frameworks!")
	if useFrameworksPos == -1 {
		return fmt.Errorf("Podfile格式不正确，未找到use_frameworks!")
	}

	// 查找use_frameworks!所在行的结尾
	lineEndPos := strings.Index(podfileContent[useFrameworksPos:], "\n")
	if lineEndPos == -1 {
		lineEndPos = len(podfileContent)
	} else {
		lineEndPos += useFrameworksPos
	}

	// 在use_frameworks!后面插入SwiftLint依赖
	swiftlintDependency := "\n  \n  # SwiftLint for code review\n  pod 'SwiftLint'"
	newContent := podfileContent[:lineEndPos] + swiftlintDependency + podfileContent[lineEndPos:]

	// 写回文件
	if err := os.WriteFile(podfilePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入Podfile失败: %w", err)
	}

	fmt.Println("  ✅ 已向Podfile添加SwiftLint依赖")
	return nil
}

// addSwiftLintToProjectYml 向project.yml添加SwiftLint构建脚本
func (s *SwiftInitializer) addSwiftLintToProjectYml() error {
	projectYmlPath := filepath.Join(s.FilePath, "project.yml")

	// 读取现有project.yml内容
	content, err := os.ReadFile(projectYmlPath)
	if err != nil {
		return fmt.Errorf("读取project.yml失败: %w", err)
	}

	projectYmlContent := string(content)

	// 检查是否已经包含SwiftLint脚本
	if strings.Contains(projectYmlContent, "SwiftLint") {
		fmt.Println("  ✅ project.yml已包含SwiftLint脚本")
		return nil
	}

	// 查找dependencies部分的结尾，在其后添加preBuildScripts
	dependenciesEndPattern := "- sdk: AVFoundation.framework"
	dependenciesPos := strings.Index(projectYmlContent, dependenciesEndPattern)
	if dependenciesPos == -1 {
		return fmt.Errorf("project.yml格式不正确，未找到dependencies部分")
	}

	// 查找该行的结尾
	lineEndPos := strings.Index(projectYmlContent[dependenciesPos:], "\n")
	if lineEndPos == -1 {
		lineEndPos = len(projectYmlContent)
	} else {
		lineEndPos += dependenciesPos
	}

	// 添加SwiftLint构建脚本，使用CocoaPods路径
	swiftlintScript := `
    preBuildScripts:
      - name: SwiftLint
        script: |
          if [ -f "${PODS_ROOT}/SwiftLint/swiftlint" ]; then
            "${PODS_ROOT}/SwiftLint/swiftlint"
          else
            echo "warning: SwiftLint not found. Make sure you have run 'pod install' and included SwiftLint in your Podfile"
          fi`

	newContent := projectYmlContent[:lineEndPos] + swiftlintScript + projectYmlContent[lineEndPos:]

	// 写回文件
	if err := os.WriteFile(projectYmlPath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("写入project.yml失败: %w", err)
	}

	fmt.Println("  ✅ 已向project.yml添加SwiftLint构建脚本（使用CocoaPods路径）")
	return nil
}

// verifySwiftLintInProjectYml 检查project.yml是否包含SwiftLint脚本
func (s *SwiftInitializer) verifySwiftLintInProjectYml() error {
	projectYmlPath := filepath.Join(s.FilePath, "project.yml")
	content, err := os.ReadFile(projectYmlPath)
	if err != nil {
		return fmt.Errorf("读取project.yml失败: %w", err)
	}

	if !strings.Contains(string(content), "SwiftLint") {
		return fmt.Errorf("project.yml中未找到SwiftLint构建脚本")
	}

	return nil
}

// ShowNextSteps Swift项目特定的后续步骤
func (s *SwiftInitializer) ShowNextSteps() {
	s.BaseInitializer.ShowNextSteps()

	fmt.Println("\n下一步：")
	if !s.NoCheck {
		fmt.Println("1. 安装依赖（包含SwiftLint）：pod install")
		fmt.Println("2. 使用 Xcode 打开 .xcworkspace 文件")
		fmt.Println("3. 设置开发者团队ID在project.yml中")
		fmt.Println("4. 运行项目，SwiftLint会自动检查代码风格")
	} else {
		fmt.Println("1. 安装依赖：pod install")
		fmt.Println("2. 使用 Xcode 打开 .xcworkspace 文件")
		fmt.Println("3. 设置开发者团队ID在project.yml中")
	}
}
