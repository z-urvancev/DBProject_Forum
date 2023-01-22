package impl

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories"
	"DBProject/pkg/errors"
)

type ForumUseCaseImpl struct {
	forumRepository  repositories.ForumRepository
	threadRepository repositories.ThreadRepository
	userRepository   repositories.UserRepository
}

func NewForumUseCase(forum repositories.ForumRepository, thread repositories.ThreadRepository, user repositories.UserRepository) *ForumUseCaseImpl {
	return &ForumUseCaseImpl{forumRepository: forum, threadRepository: thread, userRepository: user}
}

func (fuc *ForumUseCaseImpl) CreateForum(forum *models.Forum) error {
	user, err := fuc.userRepository.GetInfoAboutUser(forum.User)
	if err != nil {
		return errors.ErrUserNotFound
	}

	oldForum, err := fuc.forumRepository.GetInfoAboutForum(forum.Slug)
	if oldForum.Slug != "" {
		*forum = *oldForum
		return errors.ErrForumAlreadyExists
	}

	forum.User = user.Nickname
	err = fuc.forumRepository.CreateForum(forum)
	return err
}

func (fuc *ForumUseCaseImpl) GetInfoAboutForum(slug string) (*models.Forum, error) {
	forum, err := fuc.forumRepository.GetInfoAboutForum(slug)
	if err != nil {
		return nil, errors.ErrForumNotExist
	}
	return forum, nil
}

func (fuc *ForumUseCaseImpl) CreateForumsThread(thread *models.Thread) error {
	forum, err := fuc.forumRepository.GetInfoAboutForum(thread.Forum)
	if err != nil {
		return errors.ErrForumOrTheadNotFound
	}

	_, err = fuc.userRepository.GetInfoAboutUser(thread.Author)
	if err != nil {
		return errors.ErrForumOrTheadNotFound
	}

	currentThread, err := fuc.threadRepository.GetBySlug(thread.Slug)
	if currentThread.Slug != "" {
		*thread = *currentThread
		return errors.ErrThreadAlreadyExists
	}

	thread.Forum = forum.Slug
	err = fuc.threadRepository.CreateThread(thread)
	return err
}

func (fuc *ForumUseCaseImpl) GetForumUsers(slug string, limit int, since string, desc bool) (*models.Users, error) {
	_, err := fuc.forumRepository.GetInfoAboutForum(slug)
	if err != nil {
		return nil, errors.ErrForumNotExist
	}

	usersSlice, err := fuc.forumRepository.GetForumUsers(slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	users := new(models.Users)
	if len(*usersSlice) == 0 {
		*users = []models.User{}
	} else {
		*users = *usersSlice
	}
	return users, err
}

func (fuc *ForumUseCaseImpl) GetForumThreads(slug string, limit int, since string, desc bool) (*models.Threads, error) {
	forum, err := fuc.forumRepository.GetInfoAboutForum(slug)
	if err != nil {
		return nil, errors.ErrForumNotExist
	}

	threadsSlice, err := fuc.forumRepository.GetForumThreads(forum.Slug, limit, since, desc)
	if err != nil {
		return nil, err
	}
	threads := new(models.Threads)
	if len(*threadsSlice) == 0 {
		*threads = []models.Thread{}
	} else {
		*threads = *threadsSlice
	}
	return threads, err
}
