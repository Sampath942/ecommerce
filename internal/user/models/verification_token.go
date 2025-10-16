package models

import "time"

type VerificationToken struct {
	ID        int       `gorm:"primary key"`
	UserID    int       `gorm:"not null"`
	User      User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Token     string    `gorm:"not null;unique"`
	ExpiresAt time.Time `gorm:"not null"`
	Used      bool      `gorm:"default:false"`
}

func (VerificationToken) TableName() string {
	return "users.verification_tokens"
}
