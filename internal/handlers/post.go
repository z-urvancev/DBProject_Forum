package handlers

import (
	"DBProject/internal/models"
	"DBProject/internal/usecases"
	"DBProject/pkg/errors"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
	"net/http"
	"strconv"
)

type PostHandler struct {
	postUseCase usecases.PostUseCase
}

func NewPostHandler(postUseCase usecases.PostUseCase) *PostHandler {
	return &PostHandler{postUseCase: postUseCase}
}

func (ph *PostHandler) GetPost(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(rawId)

	related := string(ctx.QueryArgs().Peek("related"))

	postFull, err := ph.postUseCase.GetInfoAboutPost(int64(id), related)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	postFullJSON, err := jsoniter.Marshal(postFull)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, postFullJSON, http.StatusOK)
}

func (ph *PostHandler) UpdatePost(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("id").(string)
	id, err := strconv.Atoi(rawId)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	var postUpdate models.PostUpdate
	err = jsoniter.Unmarshal(ctx.Request.Body(), &postUpdate)
	if err != nil {
		errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
		return
	}

	post := &models.Post{Id: int64(id), Message: postUpdate.Message}
	err = ph.postUseCase.UpdatePost(post)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	postJSON, _ := jsoniter.Marshal(post)

	errors.CreateResponse(ctx, postJSON, http.StatusOK)
}
