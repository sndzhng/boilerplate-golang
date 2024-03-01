package usecase

import (
	"net/http"

	"github.com/sndzhng/gin-template/internal/entity"
	"github.com/sndzhng/gin-template/internal/repository"
	"github.com/sndzhng/gin-template/internal/util"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//go:generate mockgen -package=usecasemock -destination=../../mock/usecase/admin.go . Admin

type (
	Admin interface {
		Create(admin entity.Admin) error
		Delete(admin entity.Admin) error
		Get(admin entity.Admin) (entity.Admin, error)
		GetAll(adminFilter *entity.AdminFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.Admin, error)
		Initial() error
		Update(admin entity.Admin) error
	}

	adminUsecase struct {
		adminRepository repository.Admin
		roleRepository  repository.Role
	}
)

func NewAdminUsecase(
	adminRepository repository.Admin,
	roleRepository repository.Role,
) Admin {
	return &adminUsecase{
		adminRepository: adminRepository,
		roleRepository:  roleRepository,
	}
}

func (usecase *adminUsecase) Create(admin entity.Admin) error {
	if admin.Password == nil {
		return util.Error{Code: http.StatusInternalServerError, Message: "password is nil"}
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(*admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	admin.PasswordHash = &passwordHash
	err = usecase.adminRepository.Create(admin)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (usecase *adminUsecase) Delete(admin entity.Admin) error {
	err := usecase.adminRepository.Delete(admin)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (usecase *adminUsecase) Get(admin entity.Admin) (entity.Admin, error) {
	admin, err := usecase.adminRepository.Get(admin)
	if err != nil {
		switch err {
		case gorm.ErrRecordNotFound:
			return entity.Admin{}, util.Error{Code: http.StatusNotFound, Message: err.Error()}
		default:
			return entity.Admin{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}
	}

	return admin, nil
}

func (usecase *adminUsecase) GetAll(adminFilter *entity.AdminFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.Admin, error) {
	admins, err := usecase.adminRepository.GetAll(adminFilter, sortOrder, pagination)
	if err != nil {
		return []entity.Admin{}, util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	if pagination != nil {
		pagination.CalculateTotal()
	}

	return admins, nil
}

func (usecase *adminUsecase) Initial() error {
	initialID := uint64(1)
	initialRoleName := string(entity.SuperAdminRoleName)
	role := entity.Role{
		ID:   &initialID,
		Name: &initialRoleName,
	}
	_ = usecase.roleRepository.Create(role)

	initialValue := "superadmin"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(initialValue), bcrypt.DefaultCost)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}
	admin := entity.Admin{
		ID:           &initialID,
		RoleID:       &initialID,
		Username:     &initialValue,
		PasswordHash: &passwordHash,
	}
	err = usecase.adminRepository.Create(admin)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}

func (usecase *adminUsecase) Update(admin entity.Admin) error {
	if admin.Password != nil {
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(*admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
		}

		admin.PasswordHash = &passwordHash
	}

	err := usecase.adminRepository.Update(admin)
	if err != nil {
		return util.Error{Code: http.StatusInternalServerError, Message: err.Error()}
	}

	return nil
}
