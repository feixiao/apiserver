package service

import (
	"apiserver/internal/app/model"
	"apiserver/internal/app/service/user"
	"apiserver/internal/app/service/user/impl"
)

var (
	UserRepository user.Repository
)
//Init instantiate the service
func Init()  {
	UserRepository = impl.NewMysqlImpl(model.MysqlHandler)
}