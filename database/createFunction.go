//提供把数据添加进数据库的操作

package database

import (
	"gorm.io/gorm"
)

//向Users数据库添加用户信息
func UserCreate(db *gorm.DB, user *User) (err error) {
	err = db.Create(user).Error
	return
}

//添加视频信息
func VideoCreate(db *gorm.DB, video Video) (err error) {
	err = db.Create(&video).Error
	return
}

//添加赞信息
func FavoriteCreate(db *gorm.DB, fav Favorite) (err error) {
	err = db.Create(&fav).Error
	return
}

//添加评论信息
func CommentCreate(db *gorm.DB, com *Comment) (err error) {
	err = db.Create(com).Error
	return
}

//添加关注信息
func RelationCreate(db *gorm.DB, rel Relation) (err error) {
	err = db.Create(&rel).Error
	return
}
