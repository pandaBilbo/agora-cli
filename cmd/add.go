package cmd

import (
	"fmt"
	"os"

	"devex/cmd/project"

	"github.com/spf13/cobra"
)

var (
	addPath string
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "添加代码审查功能",
	Long: `为项目添加代码审查功能，包含以下功能：
  - 代码风格检查配置
  - 代码敏感信息检查工具
  - 代码审查模板
  - github CI配置

示例：
  # 为当前目录的项目添加功能
  devex add

  # 为指定路径的项目添加功能
  devex add --path /path/to/project`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("为项目添加代码审查功能\n")
		fmt.Printf("项目路径：%s\n", addPath)

		// 使用add命令专用的初始化器
		initializer, err := project.NewInitializerForAdd(addPath)
		if err != nil {
			fmt.Printf("错误：%s\n", err)
			os.Exit(1)
		}

		// 执行添加代码审查功能的步骤
		steps := []struct {
			name string
			fn   func() error
		}{
			{"复制模板文件", initializer.CopyTemplateFiles},
			{"安装Git钩子", initializer.InstallGitHooks},
		}

		for _, step := range steps {
			if err := step.fn(); err != nil {
				fmt.Printf("错误：%s失败：%s\n", step.name, err)
				os.Exit(1)
			}
		}

		initializer.ShowNextSteps()
	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// 添加命令行选项
	addCmd.Flags().StringVarP(&addPath, "path", "p", ".", "项目路径")
}
