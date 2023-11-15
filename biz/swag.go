package biz

import (
	"os"
	"os/exec"
)

// BuildApiDoc 生成接口文档
func BuildApiDoc() error {
	// 检测当前路径是否存在main.go
	cmd := exec.Command("swag", "init")
	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	filename := pwd + "/main.go"
	_, err = os.Stat(filename)
	if err != nil {
		filename = pwd + "/app/main.go"
		_, err = os.Stat(filename)
		if err != nil {
			return err
		}
		cmd.Dir = "./app"
	}
	err = cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
