package middleware

import (
	"first/internal/app/log"
	"first/internal/common"
	"first/internal/schema"

	"github.com/gin-gonic/gin"
)

func LoginUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserSession, userErr := common.SessionGet(c)
		if userErr != nil {
			log.Debug(c.Request.Context(), "session err", userErr)
			c.AbortWithStatusJSON(200,
				schema.Response{
					Code: 309999,
					Msg:  "暂无系统使用权限，请通知管理员",
				})
			return
		}
		if UserSession.ID != "" {
			log.Debug(c.Request.Context(), "mid", "请求用户写入上下文", log.H{"user_id": UserSession.ID})
			ctx := log.NewUserIDContext(c.Request.Context(), UserSession.ID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
