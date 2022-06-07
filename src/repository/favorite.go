package repository

import (
	"douyin/src/database"
	db "douyin/src/database"

	"gorm.io/gorm"
)

//添加赞信息
func FavoriteCreate(fav db.Favorite) (err error) {
	err = db.MySqlDb.Create(&fav).Error
	return
}

//删除点赞
func FavoriteDelete(userid, videoid int64) (err error) {
	err = db.MySqlDb.Where("user_id = ? AND video_id = ?", userid, videoid).Delete(&db.Favorite{}).Error
	return
}

//查询点赞是否存在
func FavoriteQueryByUserAndVideo(userid int64, videoid int64) (exist bool) {
	var nums int64
	//select count(*) from favorites where user_id = ? AND video_id = ?
	if db.MySqlDb.Table("favorites").Select("count(*)").Where("user_id = ? AND video_id = ?", userid, videoid).Scan(&nums); nums > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//使用userId查看点赞列表视频
func FavoriteQuerybyUserID(userid int64) (videolist []db.Video, nums int64) {
	db.MySqlDb.Model(&db.Video{}).Select("videos.*").Joins("inner join favorites on favorites.video_id = videos.id").Where("favorites.user_id = ? ", userid).Scan(&videolist)
	nums = int64(len(videolist))
	return
}

//添加/取消赞操作，统一更新视频点赞数、用户获赞数、用户喜欢数
func FavoriteUpdataNumbers(videoId, videoUserId, userId int64, add bool) error {
	var n int64
	if add {
		n = 1
	} else {
		n = -1
	}
	//更新视频点赞数
	if err := database.MySqlDb.Model(&database.Video{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", n)).Error; err != nil {
		return err
	}
	//更新video用户获赞数
	if err := database.MySqlDb.Model(&database.User{}).Where("id = ?", videoId).Update("total_favorite", gorm.Expr("total_favorite + ?", n)).Error; err != nil {
		return err
	}
	//更新token用户喜欢数
	if err := database.MySqlDb.Model(&database.User{}).Where("id = ?", videoId).Update("favorite_count", gorm.Expr("favorite_count + ?", n)).Error; err != nil {
		return err
	}
	return nil
}
