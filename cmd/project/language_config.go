package project

import (
	"fmt"
	"path/filepath"
)

// LanguageConfig 语言配置结构
type LanguageConfig struct {
	Name             string   // 语言名称
	DisplayName      string   // 显示名称
	TemplateCodePath string   // 模板代码路径
	ConfigPath       string   // 配置文件路径
	GlobalConfigPath string   // 全局配置路径
	RequiredCommands []string // 必需的命令行工具
}

// GetLanguageConfig 获取指定语言的配置
func GetLanguageConfig(lang string) (*LanguageConfig, error) {
	configs := getLanguageConfigs()

	if config, exists := configs[lang]; exists {
		return config, nil
	}

	return nil, fmt.Errorf("不支持的编程语言: %s。支持的语言: %s", lang, getSupportedLanguageNames())
}

// GetSupportedLanguages 获取所有支持的语言
func GetSupportedLanguages() []string {
	configs := getLanguageConfigs()
	var languages []string
	for lang := range configs {
		languages = append(languages, lang)
	}
	return languages
}

// getSupportedLanguageNames 获取支持的语言名称字符串
func getSupportedLanguageNames() string {
	languages := GetSupportedLanguages()
	if len(languages) == 0 {
		return "无"
	}
	if len(languages) == 1 {
		return languages[0]
	}

	result := ""
	for i, lang := range languages {
		if i == len(languages)-1 {
			result += " 和 " + lang
		} else if i == 0 {
			result += lang
		} else {
			result += "、" + lang
		}
	}
	return result
}

// getLanguageConfigs 获取所有语言配置（内部方法）
func getLanguageConfigs() map[string]*LanguageConfig {
	return map[string]*LanguageConfig{
		"swift": {
			Name:             "swift",
			DisplayName:      "Swift (iOS)",
			TemplateCodePath: filepath.Join("template", "swift", "code"),
			ConfigPath:       filepath.Join("template", "swift", "config"),
			GlobalConfigPath: filepath.Join("template", "global_config"),
			RequiredCommands: []string{"xcodegen", "pod"},
		},
		// Kotlin语言支持
		"kotlin": {
			Name:             "kotlin",
			DisplayName:      "Kotlin (开发中)",
			TemplateCodePath: filepath.Join("template", "kotlin", "code"),
			ConfigPath:       filepath.Join("template", "kotlin", "config"),
			GlobalConfigPath: filepath.Join("template", "global_config"),
			RequiredCommands: []string{"gradle", "java"},
		},
		// 添加新语言示例：只需要在这里添加配置，然后在factory.go中添加对应的case即可
		// "java": {
		//     Name:             "java",
		//     DisplayName:      "Java",
		//     TemplateCodePath: filepath.Join("template", "java", "code"),
		//     ConfigPath:       filepath.Join("template", "java", "config"),
		//     GlobalConfigPath: filepath.Join("template", "global_config"),
		// },
	}
}
