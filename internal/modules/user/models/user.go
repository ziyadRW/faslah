package user

import (
	"github.com/ziyadrw/faslah/internal/base"
	userEnums "github.com/ziyadrw/faslah/internal/modules/user/enums"
)

type User struct {
	base.Model

	Email        string         `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string         `gorm:"not null" json:"-"`
	Name         string         `gorm:"type:text" json:"name"`
	Role         userEnums.Type `gorm:"type:text ;default:'viewer'" json:"role"`
}
