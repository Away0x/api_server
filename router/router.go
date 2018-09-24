package router

import (
	"net/http"

	_ "github.com/Away0x/api_server/docs"
	"github.com/Away0x/api_server/handler/sd"
	"github.com/Away0x/api_server/handler/user"
	"github.com/Away0x/api_server/router/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// 加载路由
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// Middlewares
	// 在处理某些请求时可能因为程序 bug 或者其他异常情况导致程序 panic，这时候为了不影响下一次请求的调用，需要通过 gin.Recovery() 来恢复 API 服务器
	g.Use(gin.Recovery())
	// 强制浏览器不使用缓存
	g.Use(middleware.NoCache)
	// 浏览器跨域 OPTIONS 请求设置
	g.Use(middleware.Options)
	// 一些安全设置
	g.Use(middleware.Secure)
	// 其他中间件
	g.Use(mw...)

	// 404 Handler.
	g.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "The incorrect API route.")
	})

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router
	// 访问 http://localhost:8080/debug/pprof/ 可查看当前 API 服务的状态，包括 CPU 占用情况和内存使用情况
	pprof.Register(g)

	// api for authentication functionalities
	g.POST("/login", user.Login)

	// user
	u := g.Group("/v1/user")
	u.Use(middleware.AuthMiddleware()) // 组路由中间件
	{
		u.POST("", user.Create)
		u.DELETE("/:id", user.Delete)
		u.PUT("/:id", user.Update)
		u.GET("", user.List)
		u.GET("/:username", user.Get)
	}

	// API 服务器健康状态自检
	// The health check handlers
	svcd := g.Group("/sd")
	{
		svcd.GET("/health", sd.HealthCheck)
		svcd.GET("/disk", sd.DiskCheck)
		svcd.GET("/cpu", sd.CPUCheck)
		svcd.GET("/ram", sd.RAMCheck)
	}

	return g
}
