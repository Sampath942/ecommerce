package models

type User struct {
	ID              int    `gorm:"primary key"`
	Name            string `gorm:"not null"`
	Email           string `gorm:"unique;not null"`
	PhoneNumber     string `gorm:"unique"`
	Address         string
	IsAdmin         bool `gorm:"notnull"`
	IsEmailVerified bool `gorm:"default:false"`
}

func (User) TableName() string {
	return "users.users"
}
