package errors

import (
	"DBProject/internal/models"
	"errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
)

var (
	ErrForumNotExist        = errors.New("Can't find user with id ")
	ErrForumOwnerNotFound   = errors.New("Can't find user with id ")
	ErrForumAlreadyExists   = errors.New("Forum already exists ")
	ErrForumOrTheadNotFound = errors.New("Can't find user with id ")

	ErrPostNotFound              = errors.New("Can't find user with id ")
	ErrParentPostNotExist        = errors.New("Can't find user with id ")
	ErrParentPostFromOtherThread = errors.New("Can't find user with id ")

	ErrThreadAlreadyExists = errors.New("Thread already exist ")
	ErrThreadNotFound      = errors.New("Can't find user with id ")

	ErrUserAlreadyExist = errors.New("User already exist ")
	ErrUserNotFound     = errors.New("Can't find user with id ")
	ErrUserDataConflict = errors.New("Can't find user with id ")

	ErrBadInputData = errors.New("bad input data")
	ErrBadRequest   = errors.New("bad request")

	ErrNotImplemented = errors.New("not implemented")
	ErrInternal       = errors.New("internal error")
)

var errorToCode = map[error]int{
	ErrForumNotExist:        http.StatusNotFound,
	ErrForumOwnerNotFound:   http.StatusNotFound,
	ErrForumAlreadyExists:   http.StatusConflict,
	ErrForumOrTheadNotFound: http.StatusNotFound,
	ErrThreadAlreadyExists:  http.StatusConflict,

	ErrThreadAlreadyExists: http.StatusConflict,
	ErrThreadNotFound:      http.StatusNotFound,

	ErrPostNotFound:              http.StatusNotFound,
	ErrParentPostNotExist:        http.StatusNotFound,
	ErrParentPostFromOtherThread: http.StatusConflict,

	ErrUserAlreadyExist: http.StatusConflict,
	ErrUserNotFound:     http.StatusNotFound,
	ErrUserDataConflict: http.StatusConflict,

	ErrBadInputData: http.StatusBadRequest,
	ErrBadRequest:   http.StatusBadRequest,

	ErrNotImplemented: http.StatusNotImplemented,
	ErrInternal:       http.StatusInternalServerError,
}

func ConvertErrorToCode(err error) (code int) {
	code, isErrorExist := errorToCode[err]
	if !isErrorExist {
		code = http.StatusInternalServerError
	}
	return
}

func CreateResponse(ctx *fasthttp.RequestCtx, body []byte, statusCode int) {
	ctx.Response.SetBody(body)
	ctx.Response.SetStatusCode(statusCode)
	ctx.Response.Header.Set("Content-type", "application/json; charset=utf-8")
}

func CreateErrorResponse(ctx *fasthttp.RequestCtx, err error) {
	statusCode := ConvertErrorToCode(err)
	errorJSON, errMarshal := jsoniter.Marshal(models.Error{Message: err.Error()})
	if errMarshal != nil {
		statusCode = ConvertErrorToCode(ErrInternal)
		errorJSON, _ = jsoniter.Marshal(models.Error{Message: ErrInternal.Error()})
	}
	CreateResponse(ctx, errorJSON, statusCode)
}
