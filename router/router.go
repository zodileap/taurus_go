package router

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

// 路由接口
type Router interface {
	// 设置路由
	SetRouter(router *gin.RouterGroup) []gin.IRoutes
	// 设置中间件
	SetMiddleware() []gin.HandlerFunc
	SetCors() gin.HandlerFunc
}

type Response struct {
	Code    int         `json:"code"`    // 业务状态码
	Message string      `json:"message"` // 消息
	Data    interface{} `json:"data"`    // 请求数据
}

var (
	routerMu sync.RWMutex
	servers  = make(map[string]map[string]Router)
	// 用于记录服务状态
	status = make(map[string]bool)
)

//	注册一个http服务下的一个路由
//
//	如果重复传入会触发panic
//
//	示例:
//
//	type Apier struct {
//	}
//
//	func (a Apier) SetApi(router *gin.Engine) {
//		create(router)
//	}
//
//	func init() {
//		router.Register("service","code","v1", &Apier{})
//	}
//
//	func create(router *gin.Engine) {
//		router.GET("/create", func(c *gin.Context) {
//			c.String(http.StatusOK, "user")
//		})
//	}
//
// 实际url为：code/v1/create
func Register(serverName string, name string, version string, router Router) {
	routerMu.Lock()
	// 执行完毕后解锁
	defer routerMu.Unlock()
	if router == nil {
		panic(text_panic_rouet_nil)
	}
	if _, ok := servers[serverName]; !ok {
		servers[serverName] = make(map[string]Router)
		status[serverName] = false
	}
	if _, dup := servers[serverName][name]; dup {
		panic(text_panic_router_register_twice + name)
	}

	baseUrl := fmt.Sprintf("%s/%s", name, version)
	servers[serverName][baseUrl] = router
}

// 初始化路由
//
// 会遍历所有注册的API，执行SetApi方法
func InitRouter(serverName string) *gin.Engine {
	routerMu.Lock()
	// 执行完毕后解锁
	defer routerMu.Unlock()

	if _, ok := servers[serverName]; !ok {
		panic(text_panic_server_not_register)
	}
	if status[serverName] {
		panic(text_panic_server_init_twice)
	}

	var Route = gin.Default()

	for name, apier := range servers[serverName] {
		r := Route.Group(name)
		r.Use(apier.SetMiddleware()...)
		r.Use(apier.SetCors())
		apier.SetRouter(r)
	}
	status[serverName] = true

	return Route
}
