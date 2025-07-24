package controllers

import "github.com/gin-gonic/gin"

type JsonStruct struct {
	Code  int   `json:"code"`
	Msg   any   `json:"msg"`
	Data  any   `json:"data"`
	Count int64 `json:"count"`
}

func ReturnSuccess(c *gin.Context, code int, msg any, data any, count int64) {
	json := &JsonStruct{Code: code, Msg: msg, Data: data, Count: count}
	c.JSON(200, json)
}
func ReturnError(c *gin.Context, code int, msg any) {
	json := &JsonStruct{Code: code, Msg: msg}
	c.JSON(300, json)
}
