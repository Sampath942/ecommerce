package models

type Credentials struct {
	ID       int `gorm:"primary key"`
	UserID   int
	User     User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
	Password string
}

func (Credentials) TableName() string {
	return "users.credentials"
}