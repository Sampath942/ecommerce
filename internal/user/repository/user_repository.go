package repository

import (
	"errors"
	"time"

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
		IsAdmin:     false,
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

func AddVerificationToken(verificationToken models.VerificationToken, database *db.Database) error {
	return database.DB.Create(&verificationToken).Error
}

func AddUserFromGoogleID(googleID string, email string, name string, database *db.Database) error {
	user := models.User{
		Email:        email,
		GoogleID:     googleID,
		Name:         name,
		AuthProvider: "Google",
	}
	return database.DB.Create(&user).Error
}

func GetVerificationDetailsFromToken(token string, database *db.Database) (models.VerificationToken, error) {
	var verificationToken models.VerificationToken
	result := database.DB.First(&verificationToken, "token = ?", token)
	return verificationToken, result.Error
}

func GetValidVerificationDetailsFromUserID(uid int, database *db.Database) (models.VerificationToken, error) {
	var verificationToken models.VerificationToken
	result := database.DB.First(&verificationToken, "user_id = ? and used = false and expires_at > ?", uid, time.Now())
	return verificationToken, result.Error
}

func SetVerificationTokenToUsed(verificationToken models.VerificationToken, database *db.Database) error {
	return database.DB.Model(&verificationToken).Where("token = ?", verificationToken.Token).Update("used", true).Error
}

func SetUserEmailVerified(user models.User, database *db.Database) error {
	return database.DB.Model(&user).Where("id = ?", user.ID).Update("is_email_verified", true).Error
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

func GetUserByGoogleID(googleID string, database *db.Database) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "google_id = ?", googleID)
	return user, result.Error
}

func GetUserByEmail(email string, database *db.Database) (models.User, error) {
	var user models.User
	result := database.DB.First(&user, "email = ?", email)
	return user, result.Error
}

func UpdateUserGoogleID(user *models.User, googleID string, database *db.Database) error {
	err := database.DB.Model(&models.User{}).Update("google_id", googleID).Error
	if err != nil {
		return err
	}
	return database.DB.Model(&models.User{}).Update("auth_provider", "Google").Error
}

func RemoveExpiredTokens(database *db.Database) error {
	return database.DB.Where("used = ? or expires_at < ?", true, time.Now()).Delete(&models.VerificationToken{}).Error
}
