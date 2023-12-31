package entities

import (
	"github.com/LGROW101/lgrow-shop/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type IResponse interface {
	Success(code int, data any) IResponse
	Error(code int, tractId, msg string) IResponse
	Res() error
}
type Response struct {
	StatuCode int
	Data      any
	ErrorRes  *ErrorResponse
	Context   *fiber.Ctx
	IsError   bool
}
type ErrorResponse struct {
	TraceId string `json:"trace_id"`
	Msg     string `json:"message"`
}

func NewResponse(c *fiber.Ctx) IResponse {
	return &Response{
		Context: c,
	}
}
func (r *Response) Success(code int, data any) IResponse {
	r.StatuCode = code
	r.Data = data
	logger.InitLgrowLogger(r.Context, r.Data, r.StatuCode).Print().Save()
	return r
}

func (r *Response) Error(code int, tractId, msg string) IResponse {
	r.StatuCode = code
	r.ErrorRes = &ErrorResponse{
		TraceId: tractId,
		Msg:     msg,
	}
	r.IsError = true
	logger.InitLgrowLogger(r.Context, r.ErrorRes, r.StatuCode).Print().Save()
	return r
}
func (r *Response) Res() error {
	return r.Context.Status(r.StatuCode).JSON(func() any {
		if r.IsError {
			return &r.ErrorRes
		}
		return &r.Data
	}())
}
