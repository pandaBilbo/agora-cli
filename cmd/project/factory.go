package project

import (
	"fmt"
)

// NewInitializer 根据命令类型创建不同的项目初始化器
func NewInitializer(commandType, projectName, path string, noGit, noCheck bool, remote string) (Initializer, error) {
	switch commandType {
	case "init":
		// init命令：创建新项目
		return NewInitInitializer(projectName, path, noGit, noCheck, remote)
	case "add":
		// add命令：为现有项目添加代码审查功能
		return NewAddInitializer(path)
	default:
		return nil, fmt.Errorf("不支持的命令类型: %s", commandType)
	}
}

// NewInitializerForInit 专门为init命令创建初始化器（向后兼容）
func NewInitializerForInit(projectName, path string, noGit, noCheck bool, remote string) (Initializer, error) {
	return NewInitInitializer(projectName, path, noGit, noCheck, remote)
}

// NewInitializerForAdd 专门为add命令创建初始化器（向后兼容）
func NewInitializerForAdd(path string) (Initializer, error) {
	return NewAddInitializer(path)
}

// GetSupportedLanguagesFromConfig 获取支持的语言列表（使用配置系统）
func GetSupportedLanguagesFromConfig() []string {
	return GetSupportedLanguages()
}
