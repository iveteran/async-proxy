package messages

import "github.com/kataras/iris/v12"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"message"`
	Data interface{} `json:"data"`
}

func (this *Response) SetError(code int, msg string, data ...interface{}) {
	this.Code = code
	this.Msg = msg
	if data != nil {
		this.Data = data
	}
}

// 从 上下文获取，返回结果
func FromCtxGetResult(ctx iris.Context) *Response {
	rs := ctx.Values().Get("response")
	switch rs.(type) {
	case *Response:
		return rs.(*Response)
	default:
		return nil
	}
}
