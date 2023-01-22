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

type ThreadHandler struct {
	threadUseCase usecases.ThreadUseCase
}

func NewThreadHandler(threadUseCase usecases.ThreadUseCase) *ThreadHandler {
	return &ThreadHandler{threadUseCase: threadUseCase}
}

func (th *ThreadHandler) CreatePosts(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("slug_or_id").(string)

	var posts models.Posts
	err := jsoniter.Unmarshal(ctx.Request.Body(), &posts)
	if err != nil {
		errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
		return
	}

	err = th.threadUseCase.CreateNewPosts(rawId, &posts)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	postsJSON, err := jsoniter.Marshal(posts)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, postsJSON, http.StatusCreated)
}

func (th *ThreadHandler) GetThread(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("slug_or_id").(string)

	thread, err := th.threadUseCase.GetInfoAboutThread(rawId)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	threadJSON, err := jsoniter.Marshal(thread)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}
	errors.CreateResponse(ctx, threadJSON, http.StatusOK)
}

func (th *ThreadHandler) UpdateThread(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("slug_or_id").(string)

	var threadUpdate models.ThreadUpdate
	err := jsoniter.Unmarshal(ctx.Request.Body(), &threadUpdate)
	if err != nil {
		errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
		return
	}

	thread := &models.Thread{Title: threadUpdate.Title, Message: threadUpdate.Message}
	err = th.threadUseCase.UpdateThread(rawId, thread)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	threadJSON, err := jsoniter.Marshal(thread)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}
	errors.CreateResponse(ctx, threadJSON, http.StatusOK)
}

func (th *ThreadHandler) GetThreadPosts(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("slug_or_id").(string)

	since := -1
	rawSince := string(ctx.QueryArgs().Peek("since"))
	if rawSince != "" {
		var err error
		since, err = strconv.Atoi(rawSince)
		if err != nil {
			errors.CreateErrorResponse(ctx, err)
			return
		}
	}

	limitStr := string(ctx.QueryArgs().Peek("limit"))
	defaultLimit := 100
	if limitStr != "" {
		var err error
		defaultLimit, err = strconv.Atoi(limitStr)
		if err != nil {
			errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
			return
		}
	}

	descStr := string(ctx.QueryArgs().Peek("desc"))
	defaultDesc := false
	if descStr != "" {
		var err error
		defaultDesc, err = strconv.ParseBool(descStr)
		if err != nil {
			errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
			return
		}
	}

	sort := string(ctx.QueryArgs().Peek("sort"))
	if sort == "" {
		sort = "flat"
	}

	posts, err := th.threadUseCase.GetThreadPosts(rawId, defaultLimit, since, sort, defaultDesc)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	postsJSON, err := jsoniter.Marshal(posts)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, postsJSON, http.StatusOK)
}

func (th *ThreadHandler) Vote(ctx *fasthttp.RequestCtx) {
	rawId := ctx.UserValue("slug_or_id").(string)

	var vote models.Vote
	err := jsoniter.Unmarshal(ctx.Request.Body(), &vote)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	thread, err := th.threadUseCase.VoteForThread(rawId, &vote)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	threadJSON, err := jsoniter.Marshal(thread)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, threadJSON, http.StatusOK)
}
