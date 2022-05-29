//定义了删除操作的函数
package database

import (
	"gorm.io/gorm"
)

//删除点赞
func FavoriteDelete(db *gorm.DB, userid, videoid int64) (err error) {
	err = db.Where("user_id = ? AND video_id = ?", userid, videoid).Delete(&Favorite{}).Error
	return
}

//删除评论
func CommentDelete(db *gorm.DB, comment_id int64) (err error) {
	err = db.Where("id = ? ", comment_id).Delete(&Comment{}).Error
	return
}

//删除关注
func RelationDelete(db *gorm.DB, userid, touserid int64) (err error) {
	err = db.Where("user_id = ? AND to_user_id = ?", userid, touserid).Delete(&Relation{}).Error
	return
}
