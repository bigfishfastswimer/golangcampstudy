//go:build wireinject

package wire

import (
	"gitee.com/geekbang/basic-go/wire/repository"
	"gitee.com/geekbang/basic-go/wire/repository/dao"
	"github.com/google/wire"
)

//吧方法本身传进去，不是发起调用

func InitRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return &repository.UserRepository{}
}
