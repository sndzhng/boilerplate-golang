package repository

import (
	"fmt"
	"strings"

	"github.com/sndzhng/gin-template/internal/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -package=repositorymock -destination=../../mock/repository/admin.go . Admin

type (
	Admin interface {
		Create(admin entity.Admin) error
		Delete(admin entity.Admin) error
		Get(admin entity.Admin) (entity.Admin, error)
		GetAll(adminFilter *entity.AdminFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.Admin, error)
		Update(admin entity.Admin) error
	}
	adminRepository struct {
		postgresql *gorm.DB
	}
)

func NewAdminRepository(postgresql *gorm.DB) Admin {
	return &adminRepository{postgresql: postgresql}
}

func (repository *adminRepository) Create(admin entity.Admin) error {
	err := repository.postgresql.Create(&admin).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *adminRepository) Delete(admin entity.Admin) error {
	err := repository.postgresql.Delete(&admin).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *adminRepository) Get(admin entity.Admin) (entity.Admin, error) {
	err := repository.postgresql.Joins("Role").First(&admin, admin).Error
	if err != nil {
		return entity.Admin{}, err
	}

	return admin, nil
}

func (repository *adminRepository) GetAll(adminFilter *entity.AdminFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.Admin, error) {
	connection := repository.postgresql

	if adminFilter.CreateAtAfter != nil {
		connection = connection.Where("admins.create_at > ?", *adminFilter.CreateAtAfter)
	}
	if adminFilter.CreateAtBefore != nil {
		connection = connection.Where("admins.create_at < ?", *adminFilter.CreateAtBefore)
	}

	if pagination != nil {
		pagination.RecordCount = new(int64)
		err := connection.Model(&entity.Admin{}).Where(adminFilter).Count(pagination.RecordCount).Error
		if err != nil {
			return []entity.Admin{}, err
		}

		connection = connection.Limit(pagination.Limit).Offset(pagination.Offset)
	}

	if sortOrder != nil {
		connection = connection.Order(fmt.Sprintf("%s %s", sortOrder.Sort, strings.ToUpper(sortOrder.Order)))
	}

	admins := []entity.Admin{}
	err := connection.Joins("Role").Find(&admins, adminFilter).Error
	if err != nil {
		return []entity.Admin{}, err
	}

	return admins, nil
}

func (repository *adminRepository) Update(admin entity.Admin) error {
	err := repository.postgresql.Updates(&admin).Error
	if err != nil {
		return err
	}

	return nil
}
