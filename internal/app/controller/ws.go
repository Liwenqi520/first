package controller

import (
	"first/init/myredis"
	"first/internal/app/socket"
	"first/internal/schema"
	"strconv"

	"github.com/Liwenqi520/response"
	"github.com/gin-gonic/gin"
)

type Ws struct{}

// WSHandler ws连接
func (wsc Ws) WSHandlerBak(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		return
	}
	socket.Init(60)
	socket.Handler(c, id, nil)
}

func (wsc Ws) WSHandler(c *gin.Context) {
	id := c.Query("id")
	if id == "" {
		return
	}
	rdb := myredis.RedisCli()
	// 列出redis中存的所有mylist的id
	result, _ := rdb.LRange(c.Request.Context(), "mylist", 0, -1).Result()
	// 检查元素是否存在
	exists := false
	for _, val := range result {
		if val == id {
			exists = true
			break
		}
	}
	if !exists { // 不存在才添加
		rdb.RPush(c.Request.Context(), "mylist", id)
	}
	socket.Init(600)
	socket.Handler(c, id, nil)
}

// SendMsg 调用推送ws msg
func (wsc Ws) SendMsg(c *gin.Context) {
	id := c.Query("id")
	msg := c.Query("msg")
	err := socket.SendMessage(id, msg)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	response.JSON(c, nil, err)
}

func (wsc Ws) HandleMacheine(c *gin.Context) {
	var Params schema.HandleMacheineParams
	err := c.ShouldBind(&Params)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	Type := Params.Type
	cate := strconv.FormatInt(Type, 10)
	socket.SendMessage(Params.ID, cate)
	response.JSON(c, nil, err)
}
