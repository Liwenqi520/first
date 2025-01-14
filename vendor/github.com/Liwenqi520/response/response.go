package response

import (
	"github.com/Liwenqi520/errorx"
	"github.com/Liwenqi520/i18n"
	"github.com/Liwenqi520/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

type res struct {
	Code    int64
	Success bool
	Msg     string
	Data    interface{}
}

// JSON http json格式response响应
func JSON(c *gin.Context, v interface{}, err error) {
	spanCtx := logger.Start(c.Request.Context(), "ResponseJson")
	defer logger.End(spanCtx)
	lang := c.GetHeader("Accept-Language")
	r := &res{
		Code:    0,
		Success: true,
		Msg:     "",
		Data:    nil,
	}
	if err == nil {
		r.Code = 0
		r.Msg = i18n.T(lang, "success")
		if v != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    r.Code,
				"success": r.Success,
				"msg":     r.Msg,
				"data":    v,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    r.Code,
			"success": r.Success,
			"msg":     r.Msg,
		})
	} else {
		var m string
		iErr, ok := err.(*errorx.Error)
		if ok {
			r.Code = iErr.Code
			m = iErr.Msg
		} else {
			r.Code = 100000
		}
		//日志记录
		if m == "" {
			// 判断错误码类型
			// 屏蔽系统异常错误信息
			if r.Code%1000000 < 200000 {
				r.Msg = i18n.T(lang, 100000)
			} else {
				r.Msg = i18n.T(lang, int(r.Code%1000000))
			}
		} else {
			r.Msg = m
		}
		logger.Error(spanCtx, "response error", logger.Int64("code", r.Code), logger.String("msg", r.Msg))
		if v != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    r.Code,
				"success": false,
				"msg":     r.Msg,
				"data":    v,
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":    r.Code,
			"success": false,
			"msg":     r.Msg,
		})
	}
}
