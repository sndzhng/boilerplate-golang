package entity

import (
	"time"

	"gorm.io/gorm"
)

type (
	Admin struct {
		ID           *uint64        `gorm:"primaryKey" json:"id"`
		Role         *Role          `gorm:"foreignKey:RoleID" form:"-" json:"role,omitempty"`
		RoleID       *uint64        `binding:"required" form:"role_id" gorm:"index" json:"role_id"`
		CreateAt     *time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"create_at"`
		UpdateAt     *time.Time     `gorm:"default:CURRENT_TIMESTAMP" json:"update_at"`
		DeleteAt     gorm.DeletedAt `gorm:"index" json:"delete_at"`
		LastLoginAt  *time.Time     `gorm:"default:null" json:"last_login_at"`
		Username     *string        `binding:"required" form:"username" gorm:"not null;uniqueIndex" json:"username"`
		Password     *string        `binding:"required" gorm:"-" json:"password,omitempty"`
		PasswordHash *[]byte        `gorm:"not null" json:"-"`
	}
	AdminsWithNavigate struct {
		Admins     []Admin `json:"admins"`
		Pagination `json:"pagination"`
		SortOrder  `json:"sort_order"`
	}
	AdminFilter struct {
		Admin
		CreateAtAfter  *time.Time `form:"create_at_after" time_format:"2006-01-02T15:04:05" gorm:"-"`
		CreateAtBefore *time.Time `form:"create_at_before" time_format:"2006-01-02T15:04:05" gorm:"-"`
	}
)

func (admin *Admin) PreventField() {
	admin.ID = nil
	admin.LastLoginAt = nil
}
