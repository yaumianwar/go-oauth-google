package google

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/yaumianwar/go-oauth-google/model"
)

type Service struct {
	db *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return Service{db}
}

func (svc *Service) SaveUser(mUser *model.User) (error) {
	if _, err := svc.LoadUser(mUser.Email); err == nil {
		return fmt.Errorf("User already exists!")
	}
	if err := svc.db.Table("user_auth").Create(&mUser).Error; err != nil {
		return err
	}
	return nil
}

func (svc *Service) LoadUser(email string) (model.User, error) {
	var mUser model.User
	if err := svc.db.Table("user_auth").Where("email = ?", email).First(&mUser).Error; err != nil {
		return mUser, err
	}
	return mUser, nil
}
