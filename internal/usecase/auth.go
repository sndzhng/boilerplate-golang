package usecase

import (
	"net/http"
	"time"

	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/middleware"
	"github.com/sndzhng/gin-template/internal/repository"
	"github.com/sndzhng/gin-template/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockgen -package=usecasemock -destination=../../mock/usecase/auth.go . Auth

type (
	Auth interface {
		AdminLogin(login entity.Login) (entity.AccessToken, error)
		UserLogin(login entity.Login) (entity.AccessToken, error)
		UserReset(reset entity.Reset) error
	}
	authUsecase struct {
		adminRepository repository.Admin
		userRepository  repository.User
	}
)

func NewAuthUsecase(
	adminRepository repository.Admin,
	userRepository repository.User,
) Auth {
	return &authUsecase{
		adminRepository: adminRepository,
		userRepository:  userRepository,
	}
}

func (usecase *authUsecase) AdminLogin(login entity.Login) (entity.AccessToken, error) {
	admin := entity.Admin{Username: login.Username}
	admin, err := usecase.adminRepository.Get(admin)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.AccessToken{}, util.Error{Code: http.StatusNotFound, Message: err.Error()}
		} else {
			return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
	}

	err = bcrypt.CompareHashAndPassword(*admin.PasswordHash, []byte(*login.Password))
	if err != nil {
		return entity.AccessToken{}, util.Error{Code: http.StatusUnauthorized}
	}

	roles := []entity.RoleName{}
	switch *admin.RoleID {
	case 1:
		roles = append(roles, entity.SuperAdminRoleName)
	default:
		return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: "invalid role"}
	}

	accessToken, err := middleware.GenerateJWT(*admin.ID, roles...)
	if err != nil {
		return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	currentTime := time.Now()
	admin.LastLoginAt = &currentTime
	err = usecase.adminRepository.Update(admin)
	if err != nil {
		return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return entity.AccessToken{AccessToken: &accessToken}, nil
}

func (usecase *authUsecase) UserLogin(login entity.Login) (entity.AccessToken, error) {
	user := entity.User{Username: login.Username}
	user, err := usecase.userRepository.Get(user)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return entity.AccessToken{}, util.Error{Code: http.StatusNotFound, Message: err.Error()}
		} else {
			return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
	}

	err = bcrypt.CompareHashAndPassword(*user.PasswordHash, []byte(*login.Password))
	if err != nil {
		return entity.AccessToken{}, util.Error{Code: http.StatusUnauthorized}
	}

	accessToken, err := middleware.GenerateJWT(*user.ID, entity.UserRoleName)
	if err != nil {
		return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	currentTime := time.Now()
	user.LastLoginAt = &currentTime
	err = usecase.userRepository.Update(user)
	if err != nil {
		return entity.AccessToken{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return entity.AccessToken{AccessToken: &accessToken}, nil
}

func (usecase *authUsecase) UserReset(reset entity.Reset) error {
	user := entity.User{ID: reset.ID}
	user, err := usecase.userRepository.Get(user)
	if err != nil {
		return util.Error{Code: http.StatusNotFound, Message: err.Error()}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*reset.Password), bcrypt.DefaultCost)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	isResetPassword := false

	user.PasswordHash = &passwordHash
	user.IsResetPassword = &isResetPassword
	err = usecase.userRepository.Update(user)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}
