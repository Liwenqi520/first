package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func CalculateTimeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 请求进入中间件前的操作，比如记录请求开始时间
		start := time.Now()
		// 调用 Next() 会将控制权交给下一个中间件或路由处理函数
		c.Next()
		// 请求离开中间件后(即所有中间件的路由处理函数执行完毕)的操作，比如记录请求耗时
		latency := time.Since(start)
		log.Println("请求处理耗时：", latency)
	}
}