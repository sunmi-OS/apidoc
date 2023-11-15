package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/sunmi-OS/apidoc/biz"
)

// cmdUpload 同步命令
var cmdUpload = &cobra.Command{
	Use:                   "build",
	Short:                 "生成Swagger格式的API文档",
	Long:                  "生成Swagger格式的API文档",
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    false,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return fmt.Errorf("不需要任何传参数")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return biz.BuildApiDoc()
	},
}

// cmdBuild generate api doc
var cmdBuild = &cobra.Command{
	Use:                   "upload [option]",
	Short:                 "同步API文档",
	Long:                  "同步命令，支持同步到yapi和metersphere平台",
	DisableFlagsInUseLine: true,
	DisableFlagParsing:    false,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("不传参数或只传一个参数，参数只能是ms或yapi")
		}
		if len(args) == 1 {
			param := args[0]
			if param != "yapi" && param != "ms" {
				return fmt.Errorf("不传参数或只传一个参数，参数只能是ms或yapi")
			}
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		plateform := ""
		if len(args) == 1 {
			plateform = args[0]
		}
		return biz.SyncApiDoc(plateform)
	},
}

var rootCmd = &cobra.Command{
	Use: "apidoc",
}

func Execute() {
	rootCmd.AddCommand(cmdUpload)
	rootCmd.AddCommand(cmdBuild)
	rootCmd.CompletionOptions.HiddenDefaultCmd = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(2)
	}
	fmt.Println("success")
}
