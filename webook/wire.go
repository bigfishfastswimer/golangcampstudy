//go:build wireinject

package main

import (
	"gitee.com/geekbang/basic-go/webook/internal/ioc"
	"gitee.com/geekbang/basic-go/webook/internal/repository"
	"gitee.com/geekbang/basic-go/webook/internal/repository/cache"
	"gitee.com/geekbang/basic-go/webook/internal/repository/dao"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	"gitee.com/geekbang/basic-go/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebookServer() *gin.Engine {
	wire.Build(
		// initialize DB
		ioc.InitDB,
		// initialize redis
		ioc.InitRedis,
		dao.NewUserDAO,

		// initialize cache
		cache.NewCodeCache, cache.NewUserCache,

		// initialize repository
		repository.NewCodeRepository,
		repository.NewCachedUserRepository,

		// service initialization
		ioc.InitSMSService,
		service.NewCodeService, service.NewUserService,
		// handler initialization

		web.NewUserHandler, ioc.InitWebServer, ioc.InitGinMiddlewares,
	)
	return gin.Default()
}
