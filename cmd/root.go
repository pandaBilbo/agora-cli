package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devex",
	Short: "DevEx CLI 是一个帮助开发者优化开发体验的工具",
	Long: `DevEx CLI 是一个帮助开发者优化开发体验的工具。

` + `

使用示例：
  # 创建一个新的 Swift 项目
  devex init --remote https://github.com/username/myapp.git

  # 为当前目录的项目添加代码审查功能
  devex add 

  # 查看版本信息
  devex version`,
}

func init() {
	// 禁用默认的completion命令
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	// 自定义帮助模板
	rootCmd.SetUsageTemplate(`使用方法：
  {{.CommandPath}} [命令]

命令列表：
{{- range .Commands}}{{if .IsAvailableCommand}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

参数选项：
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}

查看命令详细信息请使用: "{{.CommandPath}} [命令] --help"
`)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
