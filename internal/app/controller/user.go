package controller

import (
	"first/internal/app/logic"
	"first/internal/common"
	"first/internal/schema"

	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/response"
	"github.com/gin-gonic/gin"
)

type User struct{}

// Login 登录
func (u User) Login(c *gin.Context) {
	var Params schema.LoginParams
	if err := c.ShouldBind(&Params); err != nil {
		err = errorx.New(300000, "登录参数有误")
		return
	}
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	Host := c.Request.Host
	data, err := logic.User{}.Login(c.Request.Context(), Params, Host, reqIP)
	response.JSON(c, data, err)
}

// LogOut 登出
func (u User) LogOut(c *gin.Context) {
	ctx := c.Request.Context()
	token := c.GetHeader("Authorization")
	reqIP := c.ClientIP()
	if reqIP == "::1" {
		reqIP = "127.0.0.1"
	}
	userSession, err := common.SessionGet(c)
	if err == nil {
		err = logic.User{}.DestroyToken(ctx, token, userSession, reqIP)
		response.JSON(c, nil, err)
		return
	}
	response.JSON(c, nil, err)
}