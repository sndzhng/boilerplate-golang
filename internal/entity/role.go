package entity

type (
	Role struct {
		ID   *uint64 `gorm:"primaryKey" json:"id"`
		Name *string `gorm:"not null;uniqueIndex" json:"name"`
	}
	RoleName string
)

const (
	SuperAdminRoleName RoleName = "SUPER_ADMIN"
	UserRoleName       RoleName = "USER"
)
