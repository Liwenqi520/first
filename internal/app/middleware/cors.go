package middleware

import (
	appInit "first/init"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CrossDomain() gin.HandlerFunc {
	conf := cors.DefaultConfig()
	conf.AllowOrigins = appInit.AppConfig.AllowOrigins
	conf.AllowOriginFunc = func(Origin string) bool {
		return true
	}
	conf.AddAllowMethods(
		appInit.AppConfig.AllowMethods...,
	)
	conf.AddAllowHeaders(
		appInit.AppConfig.AllowHeaders...,
	)
	conf.AllowCredentials = true
	return cors.New(conf)
}
