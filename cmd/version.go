package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// 版本信息变量（通过构建时ldflags注入）
var (
	Version    = "v0.1.0"
	BuildTime  = "unknown"
	CommitHash = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "显示版本信息",
	Long:  "显示DevEx CLI的版本、构建时间和提交哈希信息。",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("DevEx CLI %s\n", Version)
		if BuildTime != "unknown" {
			fmt.Printf("构建时间: %s\n", BuildTime)
		}
		if CommitHash != "unknown" {
			fmt.Printf("提交哈希: %s\n", CommitHash)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
