package models

import (
	"taskgo/internal/enums"
	"taskgo/pkg/utils"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	Base
	Email         string         `gorm:"uniqueIndex;not null" json:"email"`
	Password      string         `gorm:"not null" json:"-"`
	FirstName     string         `gorm:"not null" json:"first_name"`
	LastName      string         `gorm:"not null" json:"last_name"`
	Role          enums.UserRole `gorm:"type:varchar(20);not null;default:'customer'" json:"role"`
	PhoneNumber   string         `gorm:"size:20" json:"phone_number"`
	LastLoginAt   time.Time      `json:"last_login_at"`
	IsActive      bool           `gorm:"default:true" json:"is_active"`
	Notifications []Notification `gorm:"polymorphic:Notifiable;polymorphicValue:User"`
}

// BeforeCreate hooks the user model before creating a new record
func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := utils.HashPassword(u.Password)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// CheckPassword compares the provided password with the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// Implement Notifiable interface
func (u *User) GetNotifiableID() uint {
	return u.ID
}
