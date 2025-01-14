package router

import (
	"first/internal/app/controller"

	"github.com/gin-gonic/gin"
)

func OpenAiRouter(r *gin.RouterGroup) {
	r.POST("/chat-completions", controller.OpenAi{}.ChatCompletions)
	// r.POST("/completions", controller.OpenAi{}.Completions)

	r.POST("/user-recharge", controller.OpenAi{}.UserRecharge)
}
