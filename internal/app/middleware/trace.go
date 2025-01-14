package middleware

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Liwenqi520/logger"
	"github.com/gin-gonic/gin"
)

// bodyLogWriter 定义一个存储响应内容的结构体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 读取响应数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		spanCtx := logger.Start(c.Request.Context(), "trace middlerware")
		defer logger.End(spanCtx)

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw
		reqTime := time.Now()

		/*--------------*/
		//保存header和body信息
		header := c.Request.Header.Clone()
		requestBody := getRequestBody(c)
		c.Next()

		// 获取响应内容
		responseBody := blw.body.String()

		RequestUsedTime := time.Since(reqTime)
		if RequestUsedTime > time.Duration(2)*time.Second {
			logger.Error(spanCtx, "接口请求时间过长",
				logger.String("url", c.Request.URL.RequestURI()),
				logger.Int64("elapsed[ms]", RequestUsedTime.Milliseconds()),
			)
		}
		logger.Info(spanCtx, c.Request.URL.RequestURI(),
			logger.Int64("elapsed[ms]", RequestUsedTime.Milliseconds()),
			logger.Any("req-header", header),
			logger.String("req-body", requestBody),
			logger.String("response_body", responseBody),
		)
	}
}

// getRequestBody 获取请求参数
func getRequestBody(c *gin.Context) (res string) {
	switch c.Request.Method {
	case http.MethodGet:
		queryString := c.Request.URL.Query()
		qstr, _ := json.Marshal(&queryString)
		return string(qstr)

	case http.MethodPost:
		fallthrough
	case http.MethodPut:
		fallthrough
	case http.MethodPatch:
		var bodyBytes []byte // 返回body
		bodyBytes, err := c.GetRawData()
		if err != nil {
			return
		}
		// 将body还原
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

		return string(bodyBytes)

	}

	return
}
