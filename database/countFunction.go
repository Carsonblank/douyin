//提供返回点赞、评论、关注和被关注数量总数的函数
package database

import (
	"gorm.io/gorm"
)

//查询视频获得的点赞数量
func VideoFavoriteCount(db *gorm.DB, video_id int64) (nums int64) {
	//统计命中的数量
	db.Table("favorites").Select("count(*)").Where("video_id = ?", video_id).Scan(&nums)

	return
}

//查询视频的评论数量
func VideoCommentCount(db *gorm.DB, video_id int64) (nums int64) {

	db.Table("comments").Select("count(*)").Where("video_id = ?", video_id).Scan(&nums)

	return
}

//查询用户的粉丝数量
func UserFollowerCount(db *gorm.DB, user_id int64) (nums int64) {
	db.Table("relations").Select("count(*)").Where("to_user_id = ?", user_id).Scan(&nums)
	return
}

//查询用户关注的数量
func UserFollowCount(db *gorm.DB, user_id int64) (nums int64) {
	db.Table("relations").Select("count(*)").Where("user_id = ?", user_id).Scan(&nums)
	return
}

//查询用户被赞总数
func UserFavoritedCount(db *gorm.DB, user_id int64) (nums int64) {
	db.Table("favorites").Select("count(*)").Where("to_user_id = ?", user_id).Scan(&nums)
	return
}

//查询喜欢总数
func UserFavoriteCount(db *gorm.DB, user_id int64) (nums int64) {
	db.Table("favorites").Select("count(*)").Where("user_id = ?", user_id).Scan(&nums)
	return
}
