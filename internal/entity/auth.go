package entity

type (
	AccessToken struct {
		AccessToken *string `json:"access_token"`
	}
	Login struct {
		Username *string `binding:"required" json:"username"`
		Password *string `binding:"required" json:"password"`
	}
	Reset struct {
		ID       *uint64
		Password *string `binding:"required" json:"password"`
	}
)
