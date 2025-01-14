package controller

import (
	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/response"
	"first/internal/app/logic"
	"first/internal/schema"
	"fmt"
	"log"

	"first/internal/common"

	"github.com/gin-gonic/gin"
)

type OpenAi struct{}

func (ai OpenAi) Completions(c *gin.Context) {
	cli := logic.OpenAi{}.GetClinet()
	var ReqParams schema.OpenAiParams
	if err := c.ShouldBind(&ReqParams); err != nil {
		err = errorx.New(-1, err.Error())
		response.JSON(c, nil, err)
		return
	}
	uri := "/completions"
	params := map[string]interface{}{
		"model":       "text-davinci-003",
		"prompt":      ReqParams.Message,
		"max_tokens":  2048,
		"temperature": 0.9,
		"n":           1,
		"stream":      false,
	}

	res, err := cli.Post(uri, params)

	if err != nil {
		log.Fatalf("request api failed: %v", err)
	}

	fmt.Println(res.GetString("choices.0.text"))
	response.JSON(c, res.GetString("choices.0.text"), nil)
}

// UserRecharge 用户充值
func (ai OpenAi) UserRecharge(c *gin.Context) {
	Admin, err := common.SessionGet(c)
	if err != nil || Admin.ID == "" {
		err = errorx.New(-1, err.Error())
		response.JSON(c, nil, err)
		return
	}
	var ReqParams schema.OpenAiRecharge
	if err := c.ShouldBind(&ReqParams); err != nil {
		err = errorx.New(-1, err.Error())
		response.JSON(c, nil, err)
		return
	}
	message, err := logic.OpenAi{}.UserRecharge(ReqParams, Admin)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	response.JSON(c, message, nil)
}

func (ai OpenAi) ChatCompletions(c *gin.Context) {
	Admin, err := common.SessionGet(c)
	if err != nil || Admin.ID == "" {
		err = errorx.New(-1, err.Error())
		response.JSON(c, nil, err)
		return
	}
	var ReqParams schema.OpenAiParams
	if err := c.ShouldBind(&ReqParams); err != nil {
		err = errorx.New(-1, err.Error())
		response.JSON(c, nil, err)
		return
	}
	message, err := logic.OpenAi{}.ChatCompletions(ReqParams, Admin)
	if err != nil {
		response.JSON(c, nil, err)
		return
	}
	response.JSON(c, message, nil)
}
