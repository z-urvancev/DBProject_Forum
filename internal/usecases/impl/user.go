package impl

import (
	"DBProject/internal/models"
	"DBProject/internal/repositories"
	"DBProject/pkg/errors"
)

type UserUseCaseImpl struct {
	userRepository repositories.UserRepository
}

func NewUserUseCase(user repositories.UserRepository) *UserUseCaseImpl {
	return &UserUseCaseImpl{userRepository: user}
}

func (uuc *UserUseCaseImpl) CreateNewUser(user *models.User) (*models.Users, error) {
	var users *models.Users
	similarUsers, err := uuc.userRepository.GetSimilarUsers(user)
	if err != nil {
		return users, errors.ErrUserAlreadyExist
	} else if len(*similarUsers) > 0 {
		users = new(models.Users)
		*users = *similarUsers
		return users, errors.ErrUserAlreadyExist
	}
	err = uuc.userRepository.CreateUser(user)
	return users, err
}

func (uuc *UserUseCaseImpl) GetInfoAboutUser(nickname string) (*models.User, error) {
	user, err := uuc.userRepository.GetInfoAboutUser(nickname)
	if err != nil {
		return nil, errors.ErrUserNotFound
	}
	return user, nil
}

func (uuc *UserUseCaseImpl) UpdateUser(user *models.User) error {
	oldUser, err := uuc.userRepository.GetInfoAboutUser(user.Nickname)
	if oldUser.Nickname == "" {
		return errors.ErrUserNotFound
	}
	if oldUser.Fullname != user.Fullname && user.Fullname == "" {
		user.Fullname = oldUser.Fullname
	}
	if oldUser.About != user.About && user.About == "" {
		user.About = oldUser.About
	}
	if oldUser.Email != user.Email && user.Email == "" {
		user.Email = oldUser.Email
	}
	err = uuc.userRepository.UpdateUser(user)
	if err != nil {
		return errors.ErrUserDataConflict
	}
	return nil
}
