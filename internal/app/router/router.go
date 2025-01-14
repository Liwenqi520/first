package router

import (
	"first/internal/app/controller"
	"first/internal/app/middleware"

	"github.com/gin-gonic/gin"
)

func Routers(engine *gin.Engine) {
	// engine.Static("/assets", "./assets")
	// engine.LoadHTMLGlob("templates/*")

	engine.Use(middleware.CrossDomain(), middleware.TraceMiddleware())
	engine.POST("/user/login", controller.User{}.Login)
	// engine.POST("/user/send-code", controller.User{}.SendCode)
	// engine.POST("/auth", controller.User{}.Auth)
	// engine.POST("/check-token", controller.User{}.CheckToken)
	engine.POST("user/log-out", controller.User{}.LogOut)
	// engine.POST("/device/data/receive", controller.Data{}.Receive) // 数据接收接口，暂时不需要验证身份
	// engine.Any("/websocket/handle-machine", controller.Ws{}.HandleMacheine)
	// engine.Any("/websocket/ws", controller.Ws{}.WSHandler)
	// engine.Any("/websocket/ws-send-msg", controller.Ws{}.SendMsg)
	// engine.Any("/websocket/ws-test", controller.Ws{}.WsTest)

	// IndexRouter(engine.Group("index"))
	// ScanLoginRouter(engine.Group("scan-login"))

	engine.Use(middleware.CalculateTimeMiddleware(), middleware.UserMiddleware(), middleware.PermissionMiddleware())
	// UserRouter(engine.Group("user"))
	// OpenAi相关
	OpenAiRouter(engine.Group("openai"))

}
