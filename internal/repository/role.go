package repository

import (
	"github.com/sndzhng/gin-template/internal/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -package=repositorymock -destination=../../mock/repository/role.go . Role

type (
	Role interface {
		Create(role entity.Role) error
	}
	roleRepository struct {
		postgresql *gorm.DB
	}
)

func NewRoleRepository(postgresql *gorm.DB) Role {
	return &roleRepository{postgresql: postgresql}
}

func (repository *roleRepository) Create(role entity.Role) error {
	err := repository.postgresql.Create(&role).Error
	if err != nil {
		return err
	}

	return nil
}
