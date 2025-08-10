package controllers

import "github.com/gin-gonic/gin"

type JsonStruct struct {
	Code  int   `json:"code"`
	Msg   any   `json:"msg"`
	Data  any   `json:"data"`
	Count int64 `json:"count"`
}
type JsonStruct_bad struct {
	Code int `json:"code"`
	Msg  any `json:"msg"`
	Data any `json:"data"`
}

func (j JsonStruct) ReturnSuccess(c *gin.Context, code int, msg any, data any, count int64) {
	json := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	c.JSON(200, json)
}
func (j JsonStruct_bad) ReturnError(c *gin.Context, code int, msg any, data any) {
	json := &JsonStruct_bad{Code: code, Msg: msg, Data: data}
	// 根据错误类型返回相应的HTTP状态码
	httpCode := 200 // 默认200，让前端根据json中的code字段判断
	if code == 401 {
		httpCode = 401
	} else if code >= 400 {
		httpCode = code
	}
	c.JSON(httpCode, json)
}
