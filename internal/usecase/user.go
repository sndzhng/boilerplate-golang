package usecase

import (
	"net/http"

	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/repository"
	"github.com/sndzhng/gin-template/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockgen -package=usecasemock -destination=../../mock/usecase/user.go . User

type (
	User interface {
		Create(user entity.User) error
		Delete(user entity.User) error
		Get(user entity.User) (entity.User, error)
		GetAll(userFilter *entity.UserFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.User, error)
		Update(user entity.User) error
	}
	userUsecase struct {
		userRepository repository.User
	}
)

func NewUserUsecase(userRepository repository.User) User {
	return &userUsecase{userRepository: userRepository}
}

func (usecase *userUsecase) Create(user entity.User) error {
	if user.Password == nil {
		return util.Error{Code: http.StatusInternalServerError, Message: "password is nil"}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	user.PasswordHash = &passwordHash
	err = usecase.userRepository.Create(user)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (usecase *userUsecase) Delete(user entity.User) error {
	err := usecase.userRepository.Delete(user)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (usecase *userUsecase) Get(user entity.User) (entity.User, error) {
	user, err := usecase.userRepository.Get(user)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entity.User{}, util.Error{Code: http.StatusNotFound, Message: err.Error()}
		default:
			return entity.User{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
	}

	return user, nil
}

func (usecase *userUsecase) GetAll(userFilter *entity.UserFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.User, error) {
	users, err := usecase.userRepository.GetAll(userFilter, sortOrder, pagination)
	if err != nil {
		return []entity.User{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if pagination != nil {
		pagination.CalculateTotal()
	}

	return users, nil
}

func (usecase *userUsecase) Update(user entity.User) error {
	if user.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*user.Password), bcrypt.DefaultCost)
		if err != nil {
			return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
		isResetPassword := true

		user.PasswordHash = &passwordHash
		user.IsResetPassword = &isResetPassword
	}

	err := usecase.userRepository.Update(user)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}
