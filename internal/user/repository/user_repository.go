package repository

import (
	"errors"

	"github.com/Sampath942/ecommerce/db"
	"github.com/Sampath942/ecommerce/internal/user/models"
	"github.com/Sampath942/ecommerce/internal/user/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func PerformAddUserAndCredentials(addUserReq utils.AddUserRequest, database *db.Database) (models.User, error) {
	user := models.User{
		Name:        addUserReq.Name,
		Email:       addUserReq.Email,
		PhoneNumber: addUserReq.PhoneNumber,
		Address:     addUserReq.Address,
		IsAdmin:     addUserReq.IsAdmin,
	}

	err := database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&user).Error; err != nil {
			return err
		}
		if err := tx.Create(&models.Credentials{
			UserID:   user.ID,
			Password: addUserReq.Password,
		}).Error; err != nil {
			return err
		}
		return nil
	})

	return user, err
}

func GetUserById(id int, database *db.Database) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "id = ?", id)
	return user, result.Error
}

func GetUserByCredentials(email string, password string, database *db.Database) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "email = ?", email)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, result.Error
		} else {
			return models.User{}, errors.New("unable to query the database")
		}
	}
	var creds models.Credentials
	result = database.DB.First(&creds, "user_id = ?", user.ID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return models.User{}, result.Error
		} else {
			return models.User{}, errors.New("unable to query the database")
		}
	}
	err := bcrypt.CompareHashAndPassword([]byte(creds.Password), []byte(password))
	if err == nil {
		return user, nil
	}
	return models.User{}, errors.New("password doesn't match")
}
