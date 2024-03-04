package model

import (
	"time"

	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
	"github.com/achwanyusuf/carrent-lib/pkg/logger"
	"github.com/gin-gonic/gin"
)

type TransactionInfo struct {
	RequestURI    string    `json:"request_uri"`
	RequestMethod string    `json:"request_method"`
	RequestID     string    `json:"request_id"`
	Timestamp     time.Time `json:"timestamp"`
	ErrorCode     int64     `json:"error_code,omitempty"`
	Cause         string    `json:"cause,omitempty"`
}

type Response struct {
	TransactionInfo TransactionInfo `json:"transaction_info"`
	Code            int64           `json:"status_code"`
	Message         string          `json:"message,omitempty"`
	Translation     *Translation    `json:"translation,omitempty"`
}

type Translation struct {
	EN string `json:"en"`
}

type RegisterResponse struct {
	Response
	Data Account `json:"data"`
}

func (r *RegisterResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return int(r.Response.Code)
}

type LoginResponse struct {
	Response
	Auth
}

func (l *LoginResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	l.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		l.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		l.Response.Code = getErrMsg.WrappedMessage.StatusCode
		l.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		l.Response.Translation = &translation
	}

	return int(l.Response.Code)
}

type SingleAccountResponse struct {
	Response
	Data Account `json:"data"`
}

func (r *SingleAccountResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return int(r.Response.Code)
}

type EmptyResponse struct {
	Response
}

func (r *EmptyResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return int(r.Response.Code)
}

type AccountsResponse struct {
	Response
	Data       []Account  `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func (r *AccountsResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	if len(r.Data) == 0 {
		r.Data = []Account{}
	}

	return int(r.Response.Code)
}

type SingleRoleResponse struct {
	Response
	Data Role `json:"data"`
}

func (r *SingleRoleResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return int(r.Response.Code)
}

type RolesResponse struct {
	Response
	Data       []Role     `json:"data"`
	Pagination Pagination `json:"pagination"`
}

func (r *RolesResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	if len(r.Data) == 0 {
		r.Data = []Role{}
	}

	return int(r.Response.Code)
}

type SingleAccountRoleResponse struct {
	Response
	Data AccountRole `json:"data"`
}

func (r *SingleAccountRoleResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	return int(r.Response.Code)
}

type AccountRolesResponse struct {
	Response
	Data       []AccountRole `json:"data"`
	Pagination Pagination    `json:"pagination"`
}

func (r *AccountRolesResponse) Transform(ctx *gin.Context, log logger.Logger, code int, err error) int {
	r.Response = Response{
		TransactionInfo: TransactionInfo{
			RequestURI:    ctx.Request.RequestURI,
			RequestMethod: ctx.Request.Method,
			RequestID:     ctx.GetHeader("x-request-id"),
			Timestamp:     time.Now(),
		},
		Code: int64(code),
	}
	if err != nil {
		getErrMsg := errormsg.GetErrorData(err)
		r.Response.TransactionInfo.ErrorCode = getErrMsg.Code
		log.Error(ctx, errormsg.WriteErr(err))
		r.Response.Code = getErrMsg.WrappedMessage.StatusCode
		r.Response.Message = getErrMsg.WrappedMessage.Message
		translation := Translation(getErrMsg.WrappedMessage.Translation)
		r.Response.Translation = &translation
	}

	if len(r.Data) == 0 {
		r.Data = []AccountRole{}
	}

	return int(r.Response.Code)
}
