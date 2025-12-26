package models

import (
	"go-tree-hollow/pkg/utils"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email         string `gorm:"uniqueIndex;not null" json:"email"`
	Password      string `gorm:"not null" json:"-"`
	Nickname      string `gorm:"type:varchar(50)" json:"nickname"`
	AvatarURL     string `gorm:"type:varchar(1024)" json:"avatar_url"`
	BackgroundURL string `gorm:"type:varchar(1024)" json:"background_url"`
	Birthday      string `gorm:"type:varchar(20)" json:"birthday"` // YYYY-MM-DD
	Bio           string `gorm:"type:varchar(255)" json:"bio"`
	Location      string `gorm:"type:varchar(100)" json:"location"`
}

// BeforeCreate 钩子：自动加密密码
func (u *User) BeforeCreate(tx *gorm.DB) error {
	hashedPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return nil
}
