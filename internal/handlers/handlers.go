package handlers

import (
	"DBProject/internal/usecases"
)

type Handlers struct {
	ForumHandler   *ForumHandler
	PostHandler    *PostHandler
	ServiceHandler *ServiceHandler
	ThreadHandler  *ThreadHandler
	UserHandler    *UserHandler
}

func NewHandlers(UseCases *usecases.UseCases) *Handlers {
	return &Handlers{
		ForumHandler:   NewForumHandler(UseCases.ForumUseCase),
		PostHandler:    NewPostHandler(UseCases.PostUseCase),
		ServiceHandler: NewServiceHandler(UseCases.ServiceUseCase),
		ThreadHandler:  NewThreadHandler(UseCases.ThreadUseCase),
		UserHandler:    NewUserHandler(UseCases.UserUseCase),
	}
}
