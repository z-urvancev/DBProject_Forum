package handlers

import (
	"DBProject/pkg/errors"
	"fmt"
	"net/http"

	"DBProject/internal/models"
	"DBProject/internal/usecases"
	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type UserHandler struct {
	userUseCase usecases.UserUseCase
}

func NewUserHandler(userUseCase usecases.UserUseCase) *UserHandler {
	return &UserHandler{userUseCase: userUseCase}
}

func (uh *UserHandler) CreateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	var user models.User
	err := jsoniter.Unmarshal(ctx.Request.Body(), &user)
	if err != nil {
		errors.CreateErrorResponse(ctx, errors.ErrBadRequest)
		return
	}

	user.Nickname = nickname
	users, err := uh.userUseCase.CreateNewUser(&user)
	if err != nil && errors.ConvertErrorToCode(err) != http.StatusConflict {
		fmt.Println(err.Error())
		errors.CreateErrorResponse(ctx, err)
		return
	}

	if err != nil && errors.ConvertErrorToCode(err) == http.StatusConflict {
		fmt.Println(err.Error())
		usersJson, internalErr := jsoniter.Marshal(users)
		if internalErr != nil {
			errors.CreateErrorResponse(ctx, internalErr)
			return
		}
		errors.CreateResponse(ctx, usersJson, errors.ConvertErrorToCode(err))
		return
	}

	userJSON, err := jsoniter.Marshal(user)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, userJSON, http.StatusCreated)
}

func (uh *UserHandler) GetUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	user, err := uh.userUseCase.GetInfoAboutUser(nickname)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	userJSON, err := jsoniter.Marshal(user)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, userJSON, http.StatusOK)
}

func (uh *UserHandler) UpdateUser(ctx *fasthttp.RequestCtx) {
	nickname := ctx.UserValue("nickname").(string)

	userUpdate := new(models.UserUpdate)
	err := jsoniter.Unmarshal(ctx.Request.Body(), userUpdate)
	if err != nil {
		user, err := uh.userUseCase.GetInfoAboutUser(nickname)
		if err != nil {
			errors.CreateErrorResponse(ctx, err)
			return
		}

		userJSON, err := jsoniter.Marshal(user)
		if err != nil {
			errors.CreateErrorResponse(ctx, err)
			return
		}

		errors.CreateResponse(ctx, userJSON, http.StatusOK)
		return
	}

	user := &models.User{Nickname: nickname, Fullname: userUpdate.Fullname, About: userUpdate.About, Email: userUpdate.Email}

	err = uh.userUseCase.UpdateUser(user)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	userJSON, err := jsoniter.Marshal(user)
	if err != nil {
		errors.CreateErrorResponse(ctx, err)
		return
	}

	errors.CreateResponse(ctx, userJSON, http.StatusOK)
}
