package database

import (
	"douyin/common"

	"gorm.io/gorm"
)

//更新token
func UpdateToken(db *gorm.DB, user_id int64, username string) (string, error) {
	if tokenString, err := common.ReleaseToken(username); err != nil {
		return "", err
	} else {
		err = db.Model(&User{}).Where("id = ?", user_id).Update("token", tokenString).Error
		return tokenString, err
	}

}
