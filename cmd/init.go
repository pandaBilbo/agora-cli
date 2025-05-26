package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"devex/cmd/project"

	"github.com/spf13/cobra"
)

var (
	initPath    string
	initNoGit   bool
	initNoCheck bool
	initRemote  string
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化新项目",
	Long: `初始化一个新的项目，包含：
  - Git 仓库初始化
  - 代码审查配置，包含代码敏感信息检查工具，代码风格检查工具，代码审查模板
  - github CI配置
  - 项目最佳实践模板

示例：
  # 通过远程仓库初始化项目
  devex init --remote https://github.com/username/myapp.git

  # 指定路径初始化（目录会自动以仓库名命名）
  devex init --remote https://github.com/username/myapp.git --path /your/parent/dir
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// 远程仓库模式：必须指定remote
		if initRemote == "" {
			fmt.Println("错误：必须指定 --remote 远程仓库地址")
			os.Exit(1)
		}

		// 解析项目名和路径
		projectName := filepath.Base(initRemote)
		if ext := filepath.Ext(projectName); ext == ".git" {
			projectName = projectName[:len(projectName)-len(ext)]
		}
		projectPath := filepath.Join(initPath, projectName)
		if _, err := os.Stat(projectPath); !os.IsNotExist(err) {
			fmt.Printf("错误：目录已存在：%s\n", projectPath)
			os.Exit(1)
		}

		fmt.Printf("初始化项目：%s\n", projectName)

		// 使用init命令专用的初始化器
		initializer, err := project.NewInitializerForInit(projectName, projectPath, initNoGit, initNoCheck, initRemote)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		steps := []struct {
			name string
			fn   func() error
		}{
			{"克隆远程仓库", initializer.CloneRepository},
			{"复制模板文件", initializer.CopyTemplateFiles},
			// {"配置代码审查", initializer.ConfigureCodeReview},
			// {"创建项目文件", initializer.CreateProject},
			// {"初始化依赖", initializer.InitDependencies},
			{"安装 Git 钩子", initializer.InstallGitHooks},
		}

		for _, step := range steps {
			if err := step.fn(); err != nil {
				fmt.Printf("错误：%s失败：%s\n", step.name, err)
				os.RemoveAll(projectPath)
				os.Exit(1)
			}
		}

		initializer.ShowNextSteps()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().StringVarP(&initRemote, "remote", "r", "", "远程仓库地址 (必需)")
	initCmd.Flags().StringVarP(&initPath, "path", "p", ".", "项目路径")
	// initCmd.Flags().BoolVar(&initNoGit, "no-git", false, "不初始化 Git 仓库")
	initCmd.Flags().BoolVar(&initNoCheck, "no-check", false, "不添加代码审查配置")

	initCmd.MarkFlagRequired("remote")
}
