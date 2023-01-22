package impl

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories"
	"DBProject/pkg/errors"
	"strconv"
)

type ThreadUseCaseImpl struct {
	voteRepository   repositories.VoteRepository
	threadRepository repositories.ThreadRepository
	userRepository   repositories.UserRepository
	postRepository   repositories.PostRepository
}

func NewThreadUseCase(vote repositories.VoteRepository, thread repositories.ThreadRepository,
	user repositories.UserRepository, post repositories.PostRepository) *ThreadUseCaseImpl {
	return &ThreadUseCaseImpl{voteRepository: vote, threadRepository: thread, userRepository: user, postRepository: post}
}

func (tuc *ThreadUseCaseImpl) CreateNewPosts(slugOrID string, posts *models.Posts) error {
	var thread *models.Thread
	var err error
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = tuc.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = tuc.threadRepository.GetById(int64(id))
	}

	if err != nil {
		return errors.ErrThreadNotFound
	}

	if len(*posts) == 0 {
		return err
	}

	if (*posts)[0].Parent != 0 {
		var parentPost *models.Post
		parentPost, err = tuc.postRepository.GetPost((*posts)[0].Parent)
		if parentPost.Thread != thread.Id {
			return errors.ErrParentPostFromOtherThread
		}
	}
	_, err = tuc.userRepository.GetInfoAboutUser((*posts)[0].Author)
	if err != nil {
		return errors.ErrUserNotFound
	}

	err = tuc.threadRepository.CreateThreadPosts(thread, posts)
	return err
}

func (tuc *ThreadUseCaseImpl) GetInfoAboutThread(slugOrID string) (*models.Thread, error) {
	var thread *models.Thread
	var err error
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = tuc.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = tuc.threadRepository.GetById(int64(id))
	}
	if err != nil {
		return nil, errors.ErrThreadNotFound
	}
	return thread, err
}

func (tuc *ThreadUseCaseImpl) UpdateThread(slugOrID string, thread *models.Thread) error {
	id, errConv := strconv.Atoi(slugOrID)
	var currentThread *models.Thread
	var err error
	if errConv != nil {
		currentThread, err = tuc.threadRepository.GetBySlug(slugOrID)
	} else {
		currentThread, err = tuc.threadRepository.GetById(int64(id))
	}
	if err != nil {
		return errors.ErrThreadNotFound
	}
	if thread.Title != "" {
		currentThread.Title = thread.Title
	}
	if thread.Message != "" {
		currentThread.Message = thread.Message
	}
	err = tuc.threadRepository.UpdateThread(currentThread)
	if err != nil {
		return err
	}
	*thread = *currentThread
	return err
}

func (tuc *ThreadUseCaseImpl) GetThreadPosts(slugOrID string, limit, since int, sort string, desc bool) (*models.Posts, error) {
	id, errConv := strconv.Atoi(slugOrID)
	var thread *models.Thread
	var err error
	if errConv != nil {
		thread, err = tuc.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = tuc.threadRepository.GetById(int64(id))
	}

	if err != nil {
		return nil, errors.ErrThreadNotFound
	}

	postsSlice := new([]models.Post)
	switch sort {
	case "tree":
		postsSlice, err = tuc.threadRepository.GetThreadPostsTree(thread.Id, limit, since, desc)
	case "parent_tree":
		postsSlice, err = tuc.threadRepository.GetThreadPostsParentTree(thread.Id, limit, since, desc)
	default:
		postsSlice, err = tuc.threadRepository.GetThreadPostsFlat(thread.Id, limit, since, desc)
	}
	if err != nil {
		return nil, err
	}
	posts := new(models.Posts)
	if len(*postsSlice) == 0 {
		*posts = []models.Post{}
	} else {
		*posts = *postsSlice
	}
	return posts, nil
}

func (tuc *ThreadUseCaseImpl) VoteForThread(slugOrID string, vote *models.Vote) (*models.Thread, error) {
	var thread *models.Thread
	var err error
	id, errConv := strconv.Atoi(slugOrID)
	if errConv != nil {
		thread, err = tuc.threadRepository.GetBySlug(slugOrID)
	} else {
		thread, err = tuc.threadRepository.GetById(int64(id))
	}

	err = tuc.voteRepository.VoteForThread(thread.Id, vote)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	thread.Votes, err = tuc.threadRepository.GetThreadVotes(thread.Id)
	return thread, err
}
