package cmd

import (
	"fmt"

	"devex/cmd/project"
)

// Language 定义了编程语言的配置（兼容旧接口）
type Language struct {
	Name string // 语言名称
}

// getSupportedLanguages 从新架构获取支持的语言列表
func getSupportedLanguages() map[string]Language {
	result := make(map[string]Language)
	supportedLangs := project.GetSupportedLanguages()

	for _, lang := range supportedLangs {
		// 获取新架构的配置
		config, err := project.GetLanguageConfig(lang)
		if err != nil {
			continue
		}

		// 直接使用新架构的配置
		result[lang] = Language{
			Name: config.DisplayName,
		}
	}

	return result
}

// SupportedLanguages 动态获取支持的编程语言（保持向后兼容）
var SupportedLanguages = getSupportedLanguages()

// IsLanguageSupported 检查语言是否支持
func IsLanguageSupported(lang string) bool {
	supportedLangs := project.GetSupportedLanguages()
	for _, supportedLang := range supportedLangs {
		if supportedLang == lang {
			return true
		}
	}
	return false
}

// GetSupportedLanguagesText 返回支持的语言列表文本
func GetSupportedLanguagesText() string {
	text := "支持的编程语言：\n"
	languages := getSupportedLanguages()

	for key, lang := range languages {
		text += fmt.Sprintf("  - %-8s %s\n", key+":", lang.Name)
	}
	return text
}

// GetLanguageConfig 获取语言配置（保持向后兼容）
func GetLanguageConfig(lang string) (Language, error) {
	languages := getSupportedLanguages()
	if config, ok := languages[lang]; ok {
		return config, nil
	}
	return Language{}, fmt.Errorf("不支持的编程语言: %s", lang)
}
