package usecase_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sndzhng/gin-template/internal/usecase"
	repositorymock "github.com/sndzhng/gin-template/mock/repository"
)

func beforeTestAuth(test *testing.T) (
	*repositorymock.MockAdmin,
	*repositorymock.MockUser,
	usecase.Auth,
) {
	controller := gomock.NewController(test)
	defer controller.Finish()

	mockAdminRepository := repositorymock.NewMockAdmin(controller)
	mockUserRepository := repositorymock.NewMockUser(controller)
	authUsecase := usecase.NewAuthUsecase(mockAdminRepository, mockUserRepository)

	return mockAdminRepository, mockUserRepository, authUsecase
}

func TestAuthAdminLogin(test *testing.T) {
	// mockAdminRepository, _, authUsecase := beforeTestAuth(test)

	// id := uint64(1)
	// username := "username"
	// password := "password"
	// passwordHash := []byte("passwordHash")
	// login := entity.Login{
	// 	Username: &username,
	// 	Password: &password,
	// }
	// admin := entity.Admin{
	// 	Username:     login.Username,
	// 	PasswordHash: &passwordHash,
	// }

	// test.Run("Success", func(test *testing.T) {
	// 	mockAdminRepository.EXPECT().Get(gomock.Any()).Return(
	// 		entity.Admin{
	// 			ID:       &id,
	// 			RoleID:   &id,
	// 			Username: &username,
	// 		},
	// 		nil,
	// 	)

	// 	token, err := authUsecase.AdminLogin(login)
	// 	assert.NoError(test, err)
	// 	assert.NotEmpty(test, token)
	// })

	// test.Run("InternalError", func(test *testing.T) {
	// 	mockAdminRepository.EXPECT().Get(admin).Return(entity.Admin{}, errors.New("internal error"))

	// 	token, err := authUsecase.AdminLogin(login)
	// 	assert.Equal(test, http.StatusInternalServerError, err.(util.Error).Code)
	// 	assert.Empty(test, token)
	// })

	// test.Run("RecordNotFound", func(test *testing.T) {
	// 	mockAdminRepository.EXPECT().Get(admin).Return(entity.Admin{}, gorm.ErrRecordNotFound)

	// 	token, err := authUsecase.AdminLogin(login)
	// 	assert.Equal(test, http.StatusNotFound, err.(util.Error).Code)
	// 	assert.Empty(test, token)
	// })
}
