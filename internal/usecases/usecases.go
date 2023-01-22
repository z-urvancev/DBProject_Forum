package usecases

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories"
	"DBProject/internal/usecases/impl"
)

type UseCases struct {
	ForumUseCase   ForumUseCase
	PostUseCase    PostUseCase
	ServiceUseCase ServiceUseCase
	ThreadUseCase  ThreadUseCase
	UserUseCase    UserUseCase
}

func NewUseCases(repositories *repositories.Repositories) *UseCases {
	return &UseCases{
		ForumUseCase: impl.NewForumUseCase(repositories.ForumRepository, repositories.ThreadRepository, repositories.UserRepository),
		PostUseCase: impl.NewPostUseCase(repositories.ForumRepository, repositories.ThreadRepository,
			repositories.UserRepository, repositories.PostRepository),
		ServiceUseCase: impl.NewServiceUseCase(repositories.ServiceRepository),
		ThreadUseCase: impl.NewThreadUseCase(repositories.VoteRepository, repositories.ThreadRepository,
			repositories.UserRepository, repositories.PostRepository),
		UserUseCase: impl.NewUserUseCase(repositories.UserRepository),
	}
}

type ForumUseCase interface {
	CreateForum(forum *models.Forum) (err error)
	GetInfoAboutForum(slug string) (forum *models.Forum, err error)
	CreateForumsThread(thread *models.Thread) (err error)
	GetForumUsers(slug string, limit int, since string, desc bool) (users *models.Users, err error)
	GetForumThreads(slug string, limit int, since string, desc bool) (threads *models.Threads, err error)
}

type PostUseCase interface {
	GetInfoAboutPost(id int64, related string) (*models.PostFull, error)
	UpdatePost(post *models.Post) (err error)
}

type ServiceUseCase interface {
	ClearService() error
	GetService() (*models.Status, error)
}

type ThreadUseCase interface {
	CreateNewPosts(slugOrID string, posts *models.Posts) error
	GetInfoAboutThread(slugOrID string) (*models.Thread, error)
	UpdateThread(slugOrID string, thread *models.Thread) error
	GetThreadPosts(slugOrID string, limit, since int, sort string, desc bool) (*models.Posts, error)
	VoteForThread(slugOrID string, vote *models.Vote) (*models.Thread, error)
}

type UserUseCase interface {
	CreateNewUser(user *models.User) (*models.Users, error)
	GetInfoAboutUser(nickname string) (*models.User, error)
	UpdateUser(user *models.User) error
}
