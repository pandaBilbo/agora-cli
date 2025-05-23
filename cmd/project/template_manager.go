package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TemplateManager 模板管理器接口
// 通用接口，所有语言都可用
type TemplateManager interface {
	LoadTemplateCode(name string) (string, error)
	RenderTemplateCode(name string, vars map[string]string) (string, error)
	// 实用方法
	TemplateExists(name string) bool
	ListTemplates() ([]string, error)
}

// FileTemplateManager 文件模板管理器通用实现
type FileTemplateManager struct {
	TemplateCodeDir string
}

func NewFileTemplateManager(templateCodeDir string) *FileTemplateManager {
	return &FileTemplateManager{
		TemplateCodeDir: templateCodeDir,
	}
}

func (m *FileTemplateManager) LoadTemplateCode(name string) (string, error) {
	path := filepath.Join(m.TemplateCodeDir, name)
	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("读取模板文件失败: %w", err)
	}
	return string(content), nil
}

func (m *FileTemplateManager) RenderTemplateCode(name string, vars map[string]string) (string, error) {
	content, err := m.LoadTemplateCode(name)
	if err != nil {
		return "", err
	}

	// 支持多种模板变量格式，方便不同语言使用
	for k, v := range vars {
		content = strings.ReplaceAll(content, "${"+k+"}", v)
		content = strings.ReplaceAll(content, "{{"+k+"}}", v)
	}
	return content, nil
}

// TemplateExists 检查模板是否存在
func (m *FileTemplateManager) TemplateExists(name string) bool {
	path := filepath.Join(m.TemplateCodeDir, name)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// ListTemplates 列出所有可用的模板
func (m *FileTemplateManager) ListTemplates() ([]string, error) {
	var templates []string

	err := filepath.Walk(m.TemplateCodeDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			// 获取相对于模板目录的路径
			relPath, err := filepath.Rel(m.TemplateCodeDir, path)
			if err != nil {
				return err
			}
			templates = append(templates, relPath)
		}
		return nil
	})

	return templates, err
}
