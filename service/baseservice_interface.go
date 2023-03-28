package service

import (
	"github.com/dyouwan/utility/pipeline"
	"github.com/gorilla/mux"
)

// BaseServiceInterface 服务接口
type BaseServiceInterface interface {
	//GetRoutePrefix 获取路由前缀 例如:/v1/user
	GetRoutePrefix() string
	//GetRoutes 获取路由前缀下有哪些子路由 例如:/delete /add
	GetRoutes() []pipeline.Route
	//RegisterRoutes 注册路由
	RegisterRoutes(router *mux.Router)
	//Close 关闭服务
	Close()
}
