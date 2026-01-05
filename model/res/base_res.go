// Package res @Author:冯铁城 [17615007230@163.com] 2025-07-30 17:44:28
package res

import "net/http"

// Result 结果返回结构体
type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// CreateSuccess 封装创建成功结果返回
func CreateSuccess(data interface{}) *Result {
	return Success(http.StatusCreated, data)
}

// QuerySuccess 封装查询成功结果返回
func QuerySuccess(data interface{}) *Result {
	return Success(http.StatusOK, data)
}

// UpdateSuccess 封装更新成功结果返回
func UpdateSuccess(data interface{}) *Result {
	return Success(http.StatusOK, data)
}

// DeleteSuccess 封装删除成功结果返回
func DeleteSuccess(data interface{}) *Result {
	return Success(http.StatusNoContent, data)
}

// BadRequestFail 封装参数错误结果返回
func BadRequestFail(msg string) *Result {
	return Fail(http.StatusBadRequest, msg)
}

// NotFoundFail 封装未找到结果返回
func NotFoundFail(msg string) *Result {
	return Fail(http.StatusNotFound, msg)
}

// ServerFail 封装服务器错误结果返回
func ServerFail(msg string) *Result {
	return Fail(http.StatusInternalServerError, msg)
}

// UnProcessFail 封装未处理结果返回
func UnProcessFail(msg string) *Result {
	return Fail(http.StatusUnprocessableEntity, msg)
}

// UnauthorizedFail 封装未授权结果返回
func UnauthorizedFail(msg string) *Result {
	return Fail(http.StatusUnauthorized, msg)
}

// ForbiddenFail 封装权限不足结果返回
func ForbiddenFail(msg string) *Result {
	return Fail(http.StatusForbidden, msg)
}

// Success 封装成功结果返回
func Success(code int, data interface{}) *Result {
	return BuildResult(code, "success", data)
}

// Fail 封装失败结果返回
func Fail(code int, msg string) *Result {
	return BuildResult(code, msg, map[string]any{})
}

// BuildResult 封装结果返回
func BuildResult(code int, msg string, data interface{}) *Result {
	return &Result{Code: code, Msg: msg, Data: data}
}

// ToJson 结果返回转json
func (r *Result) ToJson() map[string]any {
	return map[string]any{
		"code": r.Code,
		"msg":  r.Msg,
		"data": r.Data,
	}
}
