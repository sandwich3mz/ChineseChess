package main

import (
	"chesss/config"
	"chesss/internal/logic"
	"chesss/routes"
	"context"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 导入数据
	logic.ImportData(context.Background())

	routes.RunServer()
}
