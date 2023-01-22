package repositories

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories/impl"
	"github.com/jackc/pgx"
)

type Repositories struct {
	ForumRepository   ForumRepository
	PostRepository    PostRepository
	ServiceRepository ServiceRepository
	ThreadRepository  ThreadRepository
	UserRepository    UserRepository
	VoteRepository    VoteRepository
}

func NewRepositories(db *pgx.ConnPool) *Repositories {
	return &Repositories{
		ForumRepository:   impl.NewForumRepository(db),
		PostRepository:    impl.NewPostRepository(db),
		ServiceRepository: impl.NewServiceRepository(db),
		ThreadRepository:  impl.NewThreadRepository(db),
		UserRepository:    impl.NewUserRepository(db),
		VoteRepository:    impl.NewVoteRepository(db),
	}
}

type ForumRepository interface {
	CreateForum(forum *models.Forum) (err error)
	GetInfoAboutForum(slug string) (forum *models.Forum, err error)
	GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error)
	GetForumThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error)
}

type PostRepository interface {
	GetPost(id int64) (post *models.Post, err error)
	UpdatePost(post *models.Post) (err error)
}

type ServiceRepository interface {
	ClearService() (err error)
	GetService() (status *models.Status, err error)
}

type ThreadRepository interface {
	CreateThread(thread *models.Thread) (err error)
	GetThread(slugOrId interface{}) (*models.Thread, error)
	GetThreadVotes(id int64) (votesAmount int32, err error)
	UpdateThread(thread *models.Thread) error
	CreateThreadPosts(thread *models.Thread, posts *models.Posts) error
	GetThreadPostsTree(id int64, limit, since int, desc bool) (*[]models.Post, error)
	GetThreadPostsParentTree(id int64, limit, since int, desc bool) (posts *[]models.Post, err error)
	GetThreadPostsFlat(id int64, limit, since int, desc bool) (posts *[]models.Post, err error)
	GetBySlug(slug string) (thread *models.Thread, err error)
	GetById(id int64) (thread *models.Thread, err error)
}

type UserRepository interface {
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetInfoAboutUser(nickname string) (*models.User, error)
	GetSimilarUsers(user *models.User) (*[]models.User, error)
}

type VoteRepository interface {
	VoteForThread(id int64, vote *models.Vote) error
}
