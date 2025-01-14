package middleware

import (
	"first/internal/app/log"
	"first/internal/common"
	"first/internal/schema"

	"github.com/gin-gonic/gin"
)

func UserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserSession, Err := common.SessionGet(c)
		if Err != nil {
			log.Debug(c.Request.Context(), "session err", Err)
			c.AbortWithStatusJSON(200,
				schema.Response{
					Code: 309999,
					Msg:  "用户未登录",
				})
			return
		}
		if UserSession.ID != "" {
			ctx := log.NewUserIDContext(c.Request.Context(), UserSession.ID)
			c.Request = c.Request.WithContext(ctx)
		}
		c.Next()
	}
}
