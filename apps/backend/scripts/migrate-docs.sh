#!/bin/bash

# 文档模块数据库迁移脚本

echo "执行文档模块数据库迁移..."

cd "$(dirname "$0")/.."

# 创建一个临时的 Go 程序来执行迁移
cat > /tmp/doc_migrate.go << 'EOF'
package main

import (
	"fmt"
	"metadata-platform/configs"
	document "metadata-platform/internal/module/document"
	"metadata-platform/internal/utils"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		panic(err)
	}
	
	utils.InitLogger(cfg.LogLevel, cfg.LogFilePath)
	defer utils.SyncLogger()
	
	dbManager, err := utils.NewDBManager(cfg)
	if err != nil {
		panic(err)
	}
	
	userDB := dbManager.GetUserDB()
	
	fmt.Println("Starting document module migration...")
	if err := document.MigrateDatabase(userDB); err != nil {
		panic(err)
	}
	
	fmt.Println("Seeding document module data...")
	if err := document.SeedDatabase(userDB); err != nil {
		panic(err)
	}
	
	fmt.Println("Document module migration completed successfully!")
}
EOF

# 编译并运行
go run /tmp/doc_migrate.go

# 清理
rm /tmp/doc_migrate.go

echo "迁移完成！"
