package router

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

type testRouterStruct struct {
	serverName string
	name       string
	version    string
	router     Router
}

type test struct {
}

func (t test) SetRouter(router *gin.RouterGroup) []gin.IRoutes {
	return []gin.IRoutes{
		create(router),
	}
}

func (t test) SetMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{}
}

func create(router *gin.RouterGroup) gin.IRoutes {
	return router.GET("/create", func(c *gin.Context) {
		c.String(http.StatusOK, "user")
	})
}

// TestInitRouter 测试初始化路由
func TestInitRouter(t *testing.T) {

	routers := []testRouterStruct{
		{
			serverName: "server",
			name:       "test",
			version:    "v1",
			router:     &test{},
		},
		{
			serverName: "server",
			name:       "test2",
			version:    "v1",
			router:     &test{},
		},
	}
	for _, r := range routers {
		Register(r.serverName, r.name, r.version, r.router)
	}
	InitRouter("server")

}
