package main

import (
	"first/cmd"
	appInit "first/init"
	"first/internal/defines"
)

func main() {
	// 配置初始化
	path := []string{
		"./etc",
	}
	appInit.ConfigInit(defines.AppVersion, path...)

	cmd.Execute()
}
