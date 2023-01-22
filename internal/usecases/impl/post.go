package impl

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories"
	"DBProject/pkg/errors"
	"strings"
)

type PostUsecaseImpl struct {
	forumRepository  repositories.ForumRepository
	threadRepository repositories.ThreadRepository
	userRepository   repositories.UserRepository
	postRepository   repositories.PostRepository
}

func NewPostUseCase(forum repositories.ForumRepository, thread repositories.ThreadRepository,
	user repositories.UserRepository, post repositories.PostRepository) *PostUsecaseImpl {
	return &PostUsecaseImpl{forumRepository: forum, threadRepository: thread, userRepository: user, postRepository: post}
}

func (puc *PostUsecaseImpl) GetInfoAboutPost(id int64, related string) (*models.PostFull, error) {
	fullPost := new(models.PostFull)
	var post *models.Post
	var err error
	post, err = puc.postRepository.GetPost(id)
	if err != nil {
		err = errors.ErrPostNotFound
	}
	fullPost.Post = post

	var relatedDataArr []string
	if related != "" {
		relatedDataArr = strings.Split(related, ",")
	}

	for _, data := range relatedDataArr {
		switch data {
		case "thread":
			var thread *models.Thread
			thread, err = puc.threadRepository.GetById(fullPost.Post.Thread)
			if err != nil {
				err = errors.ErrThreadNotFound
			}
			fullPost.Thread = thread
		case "user":
			var user *models.User
			user, err = puc.userRepository.GetInfoAboutUser(fullPost.Post.Author)
			if err != nil {
				err = errors.ErrUserNotFound
			}
			fullPost.Author = user
		case "forum":
			var forum *models.Forum
			forum, err = puc.forumRepository.GetInfoAboutForum(fullPost.Post.Forum)
			if err != nil {
				err = errors.ErrForumNotExist
			}
			fullPost.Forum = forum
		}
	}
	return fullPost, err
}

func (puc *PostUsecaseImpl) UpdatePost(post *models.Post) error {
	currentPost, err := puc.postRepository.GetPost(post.Id)
	if err != nil {
		return errors.ErrThreadNotFound
	}

	if post.Message != "" {
		if currentPost.Message != post.Message {
			currentPost.IsEdited = true
		}
		currentPost.Message = post.Message
		err = puc.postRepository.UpdatePost(currentPost)
		if err != nil {
			return err
		}
	}
	*post = *currentPost
	return nil
}
