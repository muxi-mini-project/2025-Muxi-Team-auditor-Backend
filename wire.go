//go:generate wire
//go:build wireinject

package main

import (
	"github.com/google/wire"
	"gorm.io/gorm"
	"muxi_auditor/client"
	"muxi_auditor/config"
	"muxi_auditor/controller"
	"muxi_auditor/ioc"
	"muxi_auditor/middleware"
	"muxi_auditor/pkg/jwt"
	"muxi_auditor/pkg/viperx"
	"muxi_auditor/repository/dao"
	"muxi_auditor/router"
	"muxi_auditor/service"
)

// wire.go

// 提供 dao.UserDAO 的 provider
func ProvideUserDAO(db *gorm.DB) dao.UserDAOInterface {
	return &dao.UserDAO{DB: db}
}

func InitWebServer(confPath string) *App {
	wire.Build(
		viperx.NewVipperSetting,
		config.NewAppConf,
		config.NewJWTConf,
		config.NewOAuthConf,
		config.NewDBConf,
		config.NewLogConf,
		config.NewCacheConf,
		config.NewPrometheusConf,
		config.NewMiddleWareConf,
		config.NewQiniuConf,
		// 初始化基础依赖
		ioc.InitDB,
		ioc.InitLogger,
		ioc.InitCache,
		ioc.InitPrometheus,
		// 初始化具体模块
		dao.NewUserDAO,
		jwt.NewRedisJWTHandler,
		service.NewAuthService,
		service.NewUserService,
		ProvideUserDAO,
		service.NewProjectService,
		service.NewItemService,
		service.NewTubeService,
		controller.NewOAuthController,
		controller.NewUserController,
		controller.NewProjectController,
		controller.NewItemController,
		controller.NewTuberController,
		client.NewOAuthClient,
		router.NewRouter,

		// 中间件
		middleware.NewAuthMiddleware,
		middleware.NewLoggerMiddleware,
		middleware.NewCorsMiddleware,
		// 应用入口
		NewApp,
	)
	return &App{}
}
