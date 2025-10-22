package models

type User struct {
	ID               int    `gorm:"primary key"`
	Name             string `gorm:"not null"`
	Email            string `gorm:"unique;not null"`
	PhoneNumber      *string `gorm:"uniqueIndex"`
	Address          string
	IsAdmin          bool   `gorm:"notnull"`
	IsEmailVerified  bool   `gorm:"default:false"`
	IsMobileVerified bool   `gorm:"default:false"`
	GoogleID         *string `gorm:"uniqueIndex"`
	AuthProvider     string
}

func (User) TableName() string {
	return "users.users"
}
