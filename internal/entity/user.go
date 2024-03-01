package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	User struct {
		ID              *uint64        `gorm:"primaryKey" json:"id"`
		Admin           *Admin         `gorm:"foreignKey:AdminID" form:"-" json:"admin,omitempty" `
		AdminID         *uint64        `binding:"required" form:"admin_id" gorm:"index;not null" json:"admin_id"`
		CreateAt        *time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"create_at"`
		UpdateAt        *time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
		DeleteAt        gorm.DeletedAt `gorm:"index" json:"delete_at"`
		LastLoginAt     *time.Time     `gorm:"default:null" json:"last_login_at"`
		Username        *string        `binding:"required" form:"username" gorm:"uniqueIndex;not null" json:"username"`
		Password        *string        `binding:"required" gorm:"-" json:"password,omitempty"`
		PasswordHash    *[]byte        `gorm:"not null" json:"-"`
		Name            *string        `binding:"required" form:"name" gorm:"not null" json:"name"`
		Phone           *string        `binding:"required" form:"phone" gorm:"uniqueIndex;not null" json:"phone"`
		IsResetPassword *bool          `form:"is_reset_password" gorm:"default:true" json:"is_reset_password"`
	}
	UsersWithNavigate struct {
		Users      []User `json:"users"`
		Pagination `json:"pagination"`
		SortOrder  `json:"sort_order"`
	}
	UserFilter struct {
		User
		CreateAtAfter  *time.Time `form:"create_at_after" time_format:"2006-01-02T15:04:05" gorm:"-"`
		CreateAtBefore *time.Time `form:"create_at_before" time_format:"2006-01-02T15:04:05" gorm:"-"`
		Search         *string    `form:"search" gorm:"-"`
	}
)

func (user *User) PreventField() {
	user.AdminID = nil
	user.ID = nil
	user.IsResetPassword = nil
	user.LastLoginAt = nil
}
