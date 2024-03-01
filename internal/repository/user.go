package repository

import (
	"fmt"
	"strings"

	"github.com/sndzhng/gin-template/internal/entity"
	"gorm.io/gorm"
)

//go:generate mockgen -package=repositorymock -destination=../../mock/repository/user.go . User

type (
	User interface {
		Create(user entity.User) error
		Delete(user entity.User) error
		Get(user entity.User) (entity.User, error)
		GetAll(userFilter *entity.UserFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.User, error)
		Update(user entity.User) error
	}
	userRepository struct {
		postgresql *gorm.DB
	}
)

func NewUserRepository(postgresql *gorm.DB) User {
	return &userRepository{postgresql: postgresql}
}

func (repository *userRepository) Create(user entity.User) error {
	err := repository.postgresql.Create(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) Delete(user entity.User) error {
	err := repository.postgresql.Delete(&user).Error
	if err != nil {
		return err
	}

	return nil
}

func (repository *userRepository) Get(user entity.User) (entity.User, error) {
	err := repository.postgresql.Joins("Admin").First(&user, user).Error
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (repository *userRepository) GetAll(userFilter *entity.UserFilter, sortOrder *entity.SortOrder, pagination *entity.Pagination) ([]entity.User, error) {
	connection := repository.postgresql

	if userFilter.CreateAtAfter != nil {
		connection = connection.Where("users.create_at > ?", *userFilter.CreateAtAfter)
	}
	if userFilter.CreateAtBefore != nil {
		connection = connection.Where("users.create_at < ?", *userFilter.CreateAtBefore)
	}
	if userFilter.Search != nil {
		search := fmt.Sprintf("%%%s%%", *userFilter.Search)
		connection = connection.Where(
			fmt.Sprintf("%s OR %s",
				"users.name LIKE ?",
				"users.username LIKE ?",
			),
			search, search,
		)

	}

	if pagination != nil {
		pagination.RecordCount = new(int64)
		err := connection.Model(&entity.User{}).Where(userFilter).Count(pagination.RecordCount).Error
		if err != nil {
			return []entity.User{}, err
		}

		connection = connection.Limit(pagination.Limit).Offset(pagination.Offset)
	}

	if sortOrder != nil {
		connection = connection.Order(fmt.Sprintf("%s %s", sortOrder.Sort, strings.ToUpper(sortOrder.Order)))
	}

	users := []entity.User{}
	err := connection.Joins("Admin").Find(&users, userFilter).Error
	if err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (repository *userRepository) Update(user entity.User) error {
	err := repository.postgresql.Updates(&user).Error
	if err != nil {
		return err
	}

	return nil
}
