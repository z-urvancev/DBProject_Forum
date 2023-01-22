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

type ForumHandler struct {
	forumUseCase usecases.ForumUseCase
}

func NewForumHandler(forumUseCase usecases.ForumUseCase) *ForumHandler {
	return &ForumHandler{forumUseCase: forumUseCase}
}

func (fh *ForumHandler) CreateForum(ctx *fasthttp.RequestCtx) {
	var forum models.Forum
	err := jsoniter.Unmarshal(ctx.Request.Body(), &forum)
	if err != nil {
		errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
		return
	}

	err = fh.forumUseCase.CreateForum(&forum)
	if err != nil && errors.ConvertErrorToCode(err) != http.StatusConflict {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	if err != nil && errors.ConvertErrorToCode(err) == http.StatusConflict {
		forumJson, _ := jsoniter.Marshal(forum)

		errors.CreateResponse(ctx, forumJson, errors.ConvertErrorToCode(err))
		return
	}

	forumJson, _ := jsoniter.Marshal(forum)
	errors.CreateResponse(ctx, forumJson, http.StatusCreated)
}

func (fh *ForumHandler) GetForum(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	forum, err := fh.forumUseCase.GetInfoAboutForum(slug)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	forumJson, _ := jsoniter.Marshal(forum)

	errors.CreateResponse(ctx, forumJson, http.StatusOK)
}

func (fh *ForumHandler) CreateThread(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	var thread models.Thread
	err := jsoniter.Unmarshal(ctx.Request.Body(), &thread)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	thread.Forum = slug

	err = fh.forumUseCase.CreateForumsThread(&thread)
	if err != nil && errors.ConvertErrorToCode(err) != http.StatusConflict {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	if err != nil && errors.ConvertErrorToCode(err) == http.StatusConflict {
		threadJson, _ := jsoniter.Marshal(thread)

		errors.CreateResponse(ctx, threadJson, errors.ConvertErrorToCode(err))
		return
	}
	threadJSON, err := jsoniter.Marshal(thread)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}
	errors.CreateResponse(ctx, threadJSON, http.StatusCreated)
}

func (fh *ForumHandler) GetForumUsers(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	limitStr := string(ctx.QueryArgs().Peek("limit"))
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
			return
		}
	}

	since := string(ctx.QueryArgs().Peek("since"))
	descStr := string(ctx.QueryArgs().Peek("desc"))
	desc := false
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
			return
		}
	}

	users, err := fh.forumUseCase.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	usersJSON, _ := jsoniter.Marshal(users)

	errors.CreateResponse(ctx, usersJSON, http.StatusOK)
}

func (fh *ForumHandler) GetForumThreads(ctx *fasthttp.RequestCtx) {
	slug := ctx.UserValue("slug").(string)

	limitStr := string(ctx.QueryArgs().Peek("limit"))
	limit := 100
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil {
			errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
			return
		}
	}

	since := string(ctx.QueryArgs().Peek("since"))
	descStr := string(ctx.QueryArgs().Peek("desc"))
	desc := false
	if descStr != "" {
		var err error
		desc, err = strconv.ParseBool(descStr)
		if err != nil {
			errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
			return
		}
	}

	threads, err := fh.forumUseCase.GetForumThreads(slug, limit, since, desc)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	threadsJSON, _ := jsoniter.Marshal(threads)

	errors.CreateResponse(ctx, threadsJSON, http.StatusOK)
}
