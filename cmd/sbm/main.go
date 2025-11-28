package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/xiaobei/singbox-manager/internal/api"
	"github.com/xiaobei/singbox-manager/internal/daemon"
	"github.com/xiaobei/singbox-manager/internal/storage"
)

var (
	version = "0.1.0"
	dataDir string
	port    int
)

func init() {
	// 获取默认数据目录
	homeDir, _ := os.UserHomeDir()
	defaultDataDir := filepath.Join(homeDir, ".singbox-manager")

	flag.StringVar(&dataDir, "data", defaultDataDir, "数据目录")
	flag.IntVar(&port, "port", 9090, "Web 服务端口")
}

func main() {
	flag.Parse()

	// 打印启动信息
	fmt.Printf("singbox-manager v%s\n", version)
	fmt.Printf("数据目录: %s\n", dataDir)
	fmt.Printf("Web 端口: %d\n", port)

	// 初始化存储
	store, err := storage.NewJSONStore(dataDir)
	if err != nil {
		log.Fatalf("初始化存储失败: %v", err)
	}

	// 获取设置
	settings := store.GetSettings()

	// 初始化进程管理器
	configPath := filepath.Join(dataDir, "generated", "config.json")
	processManager := daemon.NewProcessManager(settings.SingBoxPath, configPath)

	// 初始化 launchd 管理器
	launchdManager, err := daemon.NewLaunchdManager()
	if err != nil {
		log.Printf("初始化 launchd 管理器失败: %v", err)
	}

	// 创建 API 服务器
	server := api.NewServer(store, processManager, launchdManager)

	// 启动定时任务调度器
	server.StartScheduler()

	// 启动服务
	addr := fmt.Sprintf(":%d", port)
	fmt.Printf("启动 Web 服务: http://127.0.0.1%s\n", addr)

	if err := server.Run(addr); err != nil {
		log.Fatalf("启动服务失败: %v", err)
	}
}
