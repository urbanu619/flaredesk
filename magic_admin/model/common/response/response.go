package response

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"reflect"
	"time"
)

var emptyObj = struct{}{}
var emptyArray = []struct{}{}

// 响应码

type ResponseCode struct {
	Code int32  `json:"code"`    // 响应码
	Msg  string `json:"message"` // 响应消息
}

var (
	ResponseCodeInsufficientAuthority = &ResponseCode{Code: 301, Msg: "权限不足"}
	ResponseCodeSignError             = &ResponseCode{Code: 302, Msg: "验签失败"}
	ResponseCodeMissAuthToken         = &ResponseCode{Code: 401, Msg: "Authorization不存在"}
	ResponseCodeTokenInvalid          = &ResponseCode{Code: 401, Msg: "Token过期失效"}
)

var (
	ResponseCodeSuccess           = &ResponseCode{Code: 200, Msg: "操作成功"}
	ResponseCodeFailure           = &ResponseCode{Code: 500, Msg: "操作失败"}
	ResponseCodeParamError        = &ResponseCode{Code: 502, Msg: "参数错误"}
	ResponseCodeFrequentOperation = &ResponseCode{Code: 503, Msg: "请求频繁"} // 操作频繁
)

func (rc *ResponseCode) Error() string {
	return rc.Msg
}

func (rc *ResponseCode) ReplaceMsg(msg string) *ResponseCode {
	rc.Msg = msg
	return rc
}

// 响应结构

type Response struct {
	Code int32  `json:"code"`           // 响应码
	Msg  string `json:"msg"`            // 响应消息
	Time int64  `json:"time"`           // 服务器时间
	Data any    `json:"data,omitempty"` // 响应数据
}

func (r *Response) MarshalJSON() ([]byte, error) {
	var resp struct {
		Code int32  `json:"code"` // 返回码
		Msg  string `json:"msg"`  // 返回消息
		Time int64  `json:"time"` // 服务器时间
		Data any    `json:"data"` // 返回数据
	}

	resp.Code = r.Code
	resp.Msg = r.Msg
	resp.Time = r.Time

	if r.Data == nil {
		resp.Data = emptyObj
	} else {
		value := reflect.ValueOf(r.Data)
		if value.Kind() == reflect.Slice {
			if value.Len() == 0 {
				resp.Data = emptyArray
			} else {
				resp.Data = r.Data
			}
		} else {
			resp.Data = r.Data
		}
	}

	return json.Marshal(resp)
}

func dataToResponse(data any) *Response {
	res := &Response{
		Code: ResponseCodeSuccess.Code,
		Data: data,
		Msg:  ResponseCodeSuccess.Msg,
		Time: currentTimestamp(),
	}

	return res
}

func currentTimestamp() int64 {
	return time.Now().Unix()
}

func success() *Response {
	res := &Response{
		Code: ResponseCodeSuccess.Code,
		Msg:  ResponseCodeSuccess.Msg,
		Time: currentTimestamp(),
	}
	return res
}

func failureWithMsg(msg string) *Response {
	res := &Response{
		Code: ResponseCodeFailure.Code,
		Msg:  msg,
	}

	return res
}

func ErrorObjByCode(rc *ResponseCode) *Response {
	res := &Response{
		Code: rc.Code,
		Msg:  rc.Msg,
		Time: currentTimestamp(),
	}
	return res
}

func errorWithCode(rc *ResponseCode) *Response {
	res := &Response{
		Code: rc.Code,
		Msg:  rc.Msg,
		Time: currentTimestamp(),
	}

	return res
}

// 分页列表

type PageList struct {
	PageNum  int64 `json:"pageNum,string"`  // 页码
	PageSize int64 `json:"pageSize,string"` // 分页大小
	Total    int64 `json:"total,string"`    // 总记录数
	List     any   `json:"list"`            // 数据列表
	Menu     any   `json:"menu"`            // 类型列表
}

// 用于直接返回

func respOkData(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, dataToResponse(resp))
}

func respOk(c *gin.Context) {
	c.JSON(http.StatusOK, success())
}

// 用于code多语言

func respFailureCode(c *gin.Context, respObj *ResponseCode) {
	c.JSON(http.StatusOK, errorWithCode(respObj))
}

func respFailureMsg(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, failureWithMsg(msg))
}

func Resp(c *gin.Context, args ...interface{}) {
	var data interface{}
	if len(args) > 0 {
		data = args[0]
	} else {
		respOk(c)
		return
	}
	switch data.(type) {
	case string:
		respFailureMsg(c, data.(string))
	case *ResponseCode:
		respFailureCode(c, data.(*ResponseCode))
	case error:
		respFailureMsg(c, data.(error).Error())
	default:
		respOkData(c, data)
	}
}
