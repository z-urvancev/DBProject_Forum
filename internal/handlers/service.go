package handlers

import (
	"DBProject/internal/usecases"
	"DBProject/pkg/errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ServiceHandler struct {
	serviceUseCase usecases.ServiceUseCase
}

func NewServiceHandler(serviceUseCase usecases.ServiceUseCase) *ServiceHandler {
	return &ServiceHandler{serviceUseCase: serviceUseCase}
}

func (sh *ServiceHandler) Clear(ctx *fasthttp.RequestCtx) {
	err := sh.serviceUseCase.ClearService()
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}
	ctx.Response.SetStatusCode(http.StatusOK)
}

func (sh *ServiceHandler) GetStatus(ctx *fasthttp.RequestCtx) {
	status, err := sh.serviceUseCase.GetService()
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	statusJSON, err := jsoniter.Marshal(status)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}
	errors.CreateResponse(ctx, statusJSON, http.StatusOK)
}
