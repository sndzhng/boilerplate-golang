package usecase_test

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/usecase"
	"github.com/sndzhng/gin-template/internal/util"
	repositorymock "github.com/sndzhng/gin-template/mock/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func beforeTestUser(test *testing.T) (
	*repositorymock.MockUser,
	usecase.User,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockUserRepository := repositorymock.NewMockUser(controller)
	userUsecase := usecase.NewUserUsecase(mockUserRepository)

	return mockUserRepository, userUsecase
}

func TestUserCreate(test *testing.T) {
	mockUserRepository, userUsecase := beforeTestUser(test)

	id := uint64(1)
	username := "username"
	password := "password"
	name := "name"
	user := entity.User{
		AdminID:  &id,
		Username: &username,
		Password: &password,
		Name:     &name,
	}

	test.Run("Success", func(test *testing.T) {
		mockUserRepository.EXPECT().Create(gomock.Any()).Return(nil)

		err := userUsecase.Create(user)
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserRepository.EXPECT().Create(gomock.Any()).Return(errors.New("internal error"))

		err := userUsecase.Create(user)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})

	test.Run("PasswordIsNilError", func(test *testing.T) {
		user.Password = nil

		err := userUsecase.Create(entity.User{})
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}

func TestUserDelete(test *testing.T) {
	mockUserRepository, userUsecase := beforeTestUser(test)

	id := uint64(1)
	user := entity.User{
		ID: &id,
	}

	test.Run("Success", func(test *testing.T) {
		mockUserRepository.EXPECT().Delete(user).Return(nil)

		err := userUsecase.Delete(user)
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserRepository.EXPECT().Delete(user).Return(errors.New("internal error"))

		err := userUsecase.Delete(user)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}

func TestUserGet(test *testing.T) {
	mockUserRepository, userUsecase := beforeTestUser(test)

	id := uint64(1)
	username := "username"
	user := entity.User{
		ID: &id,
	}

	test.Run("Success", func(test *testing.T) {
		mockUserRepository.EXPECT().Get(user).Return(
			entity.User{
				ID:       &id,
				Username: &username,
			},
			nil,
		)

		result, err := userUsecase.Get(user)
		assert.NoError(test, err)
		assert.Equal(test, *result.ID, *user.ID)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserRepository.EXPECT().Get(user).Return(entity.User{}, errors.New("internal error"))

		result, err := userUsecase.Get(user)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
		assert.Equal(test, result, entity.User{})
	})

	test.Run("RecordNotFound", func(test *testing.T) {
		mockUserRepository.EXPECT().Get(user).Return(entity.User{}, gorm.ErrRecordNotFound)

		result, err := userUsecase.Get(user)
		assert.Equal(test, http.StatusNotFound, err.(util.Error).Code)
		assert.Equal(test, result, entity.User{})
	})
}

func TestUserGetAll(test *testing.T) {
	mockUserRepository, userUsecase := beforeTestUser(test)

	id := uint64(1)
	username := "username"
	users := []entity.User{
		{
			ID:       &id,
			Username: &username,
		},
	}

	test.Run("Success", func(test *testing.T) {
		userFilter := entity.UserFilter{
			User:           users[0],
			CreateAtAfter:  &time.Time{},
			CreateAtBefore: &time.Time{},
		}
		sortOrder := entity.InitialSortOrder()
		pagination := entity.Pagination{
			Limit:  1,
			Offset: 0,
		}

		mockUserRepository.EXPECT().GetAll(&userFilter, &sortOrder, &pagination).Return(users, nil)

		result, err := userUsecase.GetAll(&userFilter, &sortOrder, &pagination)
		assert.NoError(test, err)
		assert.Len(test, result, len(users))
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserRepository.EXPECT().GetAll(nil, nil, nil).Return([]entity.User{}, errors.New("internal error"))

		result, err := userUsecase.GetAll(nil, nil, nil)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
		assert.Len(test, result, 0)
	})
}

func TestUserUpdate(test *testing.T) {
	mockUserRepository, userUsecase := beforeTestUser(test)

	id := uint64(1)
	password := "password"
	username := "username"
	name := "name"
	phone := "+66987654321"
	user := entity.User{
		ID:              &id,
		Username:        &username,
		Password:        &password,
		Name:            &name,
		Phone:           &phone,
		IsResetPassword: new(bool),
	}

	test.Run("Success", func(test *testing.T) {
		mockUserRepository.EXPECT().Update(gomock.Any()).Return(nil)

		err := userUsecase.Update(user)
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockUserRepository.EXPECT().Update(gomock.Any()).Return(errors.New("internal error"))

		err := userUsecase.Update(user)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}
