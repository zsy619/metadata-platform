package sso

import (
	"fmt"
	"os"

	"github.com/cloudwego/hertz/pkg/app/server"
	"gorm.io/gorm"

	"metadata-platform/internal/module/sso/api"
	"metadata-platform/internal/module/sso/repository"
	"metadata-platform/internal/module/sso/service"
)

// RegisterRoutes 注册SSO模块路由
func RegisterRoutes(r *server.Hertz, db *gorm.DB) {
	fmt.Fprintln(os.Stderr, ">>> Initializing SSO Routes...")

	fmt.Fprintln(os.Stderr, ">>> Creating repositories...")
	repos := repository.NewRepositories(db)
	fmt.Fprintln(os.Stderr, ">>> Repositories created successfully")

	fmt.Fprintln(os.Stderr, ">>> Creating services...")
	services := service.NewServices(repos, db)
	fmt.Fprintln(os.Stderr, ">>> Services created successfully")

	fmt.Fprintln(os.Stderr, ">>> Creating handlers...")
	handlers := api.NewSsoHandler(services)
	fmt.Fprintln(os.Stderr, ">>> Handlers created successfully")

	fmt.Fprintln(os.Stderr, ">>> Registering routes...")
	handlers.RegisterRoutes(r)
	fmt.Fprintln(os.Stderr, ">>> Routes registered successfully")
}
