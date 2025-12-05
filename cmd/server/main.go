package main

import (
    "log"
    "go-tree-hollow/configs"
    "go-tree-hollow/internal/server"
)

func main() {
    // 加载配置
    config := configs.LoadConfig()

    // 创建服务器
    srv, err := server.NewServer(config)
    if err != nil {
        log.Fatal("Failed to create server:", err)
    }

    // 启动服务器
    if err := srv.Start(); err != nil {
        log.Fatal("Failed to start server:", err)
    }
}