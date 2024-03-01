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

func beforeTestAdmin(test *testing.T) (
	*repositorymock.MockAdmin,
	*repositorymock.MockRole,
	usecase.Admin,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockAdminRepository := repositorymock.NewMockAdmin(controller)
	mockRoleRepository := repositorymock.NewMockRole(controller)
	adminUsecase := usecase.NewAdminUsecase(mockAdminRepository, mockRoleRepository)

	return mockAdminRepository, mockRoleRepository, adminUsecase
}

func TestAdminCreate(test *testing.T) {
	mockAdminRepository, _, adminUsecase := beforeTestAdmin(test)

	roleID := uint64(1)
	password := "password"
	username := "username"
	admin := entity.Admin{
		RoleID:   &roleID,
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		mockAdminRepository.EXPECT().Create(gomock.Any()).Return(nil)

		err := adminUsecase.Create(admin)
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminRepository.EXPECT().Create(gomock.Any()).Return(errors.New("internal error"))

		err := adminUsecase.Create(admin)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})

	test.Run("PasswordIsNil", func(test *testing.T) {
		admin.Password = nil

		err := adminUsecase.Create(admin)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}

func TestAdminDelete(test *testing.T) {
	mockAdminRepository, _, adminUsecase := beforeTestAdmin(test)

	id := uint64(1)
	admin := entity.Admin{
		ID: &id,
	}

	test.Run("Success", func(test *testing.T) {
		mockAdminRepository.EXPECT().Delete(admin).Return(nil)

		err := adminUsecase.Delete(admin)
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminRepository.EXPECT().Delete(admin).Return(errors.New("internal error"))

		err := adminUsecase.Delete(admin)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}

func TestAdminGet(test *testing.T) {
	mockAdminRepository, _, adminUsecase := beforeTestAdmin(test)

	id := uint64(1)
	username := "username"
	admin := entity.Admin{
		ID: &id,
	}

	test.Run("Success", func(test *testing.T) {
		mockAdminRepository.EXPECT().Get(admin).Return(
			entity.Admin{
				ID:       &id,
				RoleID:   &id,
				Username: &username,
			},
			nil,
		)

		result, err := adminUsecase.Get(admin)
		assert.NoError(test, err)
		assert.Equal(test, *result.ID, *admin.ID)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminRepository.EXPECT().Get(admin).Return(entity.Admin{}, errors.New("internal error"))

		result, err := adminUsecase.Get(admin)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
		assert.Equal(test, result, entity.Admin{})
	})

	test.Run("RecordNotFound", func(test *testing.T) {
		mockAdminRepository.EXPECT().Get(admin).Return(entity.Admin{}, gorm.ErrRecordNotFound)

		result, err := adminUsecase.Get(admin)
		assert.Equal(test, http.StatusNotFound, err.(util.Error).Code)
		assert.Equal(test, result, entity.Admin{})
	})
}

func TestAdminGetAll(test *testing.T) {
	mockAdminRepository, _, adminUsecase := beforeTestAdmin(test)

	id := uint64(1)
	username := "username"
	admins := []entity.Admin{
		{
			ID:       &id,
			RoleID:   &id,
			Username: &username,
		},
	}

	test.Run("Success", func(test *testing.T) {
		adminFilter := entity.AdminFilter{
			Admin:          admins[0],
			CreateAtAfter:  &time.Time{},
			CreateAtBefore: &time.Time{},
		}
		sortOrder := entity.InitialSortOrder()
		pagination := entity.Pagination{
			Limit:  1,
			Offset: 0,
		}

		mockAdminRepository.EXPECT().GetAll(&adminFilter, &sortOrder, &pagination).Return(admins, nil)

		result, err := adminUsecase.GetAll(&adminFilter, &sortOrder, &pagination)
		assert.NoError(test, err)
		assert.Len(test, result, len(admins))
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminRepository.EXPECT().GetAll(nil, nil, nil).Return([]entity.Admin{}, errors.New("internal error"))

		result, err := adminUsecase.GetAll(nil, nil, nil)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
		assert.Len(test, result, 0)
	})
}

func TestAdminInitial(test *testing.T) {
	mockAdminRepository, mockRoleRepository, adminUsecase := beforeTestAdmin(test)

	test.Run("Success", func(test *testing.T) {
		mockRoleRepository.EXPECT().Create(gomock.Any()).Return(nil)
		mockAdminRepository.EXPECT().Create(gomock.Any()).Return(nil)

		err := adminUsecase.Initial()
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockRoleRepository.EXPECT().Create(gomock.Any()).Return(nil)
		mockAdminRepository.EXPECT().Create(gomock.Any()).Return(errors.New("internal error"))

		err := adminUsecase.Initial()
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}

func TestAdminUpdate(test *testing.T) {
	mockAdminRepository, _, adminUsecase := beforeTestAdmin(test)

	id := uint64(1)
	password := "password"
	username := "username"
	admin := entity.Admin{
		ID:       &id,
		RoleID:   &id,
		Username: &username,
		Password: &password,
	}

	test.Run("Success", func(test *testing.T) {
		mockAdminRepository.EXPECT().Update(gomock.Any()).Return(nil)

		err := adminUsecase.Update(admin)
		assert.NoError(test, err)
	})

	test.Run("InternalError", func(test *testing.T) {
		mockAdminRepository.EXPECT().Update(gomock.Any()).Return(errors.New("internal error"))

		err := adminUsecase.Update(admin)
		assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	})
}
