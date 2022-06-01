//提供数据库查找操作

package database

import (
	"douyin/config"

	"gorm.io/gorm"
)

//根据videoId查看video是否存在
func VideoQueryByID(db *gorm.DB, video_id int64) (exist bool, user_id int64) {
	var v []Video
	db.Where("id = ? ", video_id).Find(&v)
	if len(v) > 0 {
		exist = true
		user_id = v[0].UserId
	} else {
		exist = false
		user_id = 0
	}
	return

}

//返回视频上传时间不大于给定时间的最多nums个视频的列表
func VideoQueryByLastTime(db *gorm.DB, lastest_time int64) []Video {
	videolist_db := make([]Video, 0, config.MustVideosNums)
	//降序排列
	db.Limit(config.MustVideosNums).Order("created_at DESC").Where("created_at <= ?", lastest_time).Find(&videolist_db)

	return videolist_db
}

//使用UserID查询视频列表
func VideoQueryByUserID(db *gorm.DB, userID int64) (videolist []Video) {
	db.Where("user_id = ? ", userID).Find(&videolist)
	return
}

//使用主键id查d询单个用户
func UserQueryByID(db *gorm.DB, id int64) (user []User, exist bool) {
	if db.Limit(1).Find(&user, id); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//使用Name查询单个用户
func UserQueryByName(db *gorm.DB, username string) (user []User, exist bool) {
	if db.Limit(1).Where("name = ?", username).Find(&user); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//验证登录信息是否正确
func UserValid(db *gorm.DB, username, password string) (user []User, exist bool) {
	if db.Limit(1).Where("name = ? AND password = ?", username, password).Find(&user); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

/*
//根据token查找user
func UserQueryByToken(db *gorm.DB, token string) (user []User, exist bool) {
	if db.Limit(1).Where("token = ?", token).Find(&user); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}
*/
//根据video_id查找user
func UserQueryByVideoID(db *gorm.DB, video_id int64) (user_id int64) {
	db.Table("videos").Select("user_id").Where("id = ? ", video_id).Scan(&user_id)

	return
}

//查询是否点赞
func FavoriteQueryByUserAndVideo(db *gorm.DB, userid int64, videoid int64) (exist bool) {
	var nums int64
	//select count(*) from favorites where user_id = ? AND video_id = ?
	if db.Table("favorites").Select("count(*)").Where("user_id = ? AND video_id = ?", userid, videoid).Scan(&nums); nums > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//查询是否关注
func FavoriteQueryByUserAndUser(db *gorm.DB, userid, touserid int64) (exist bool) {
	var nums int64
	if db.Table("relations").Select("count(*)").Where("user_id = ? AND to_user_id = ?", userid, touserid).Scan(&nums); nums > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//使用videoId查询评论列表
func CommentQuerybyVideoID(db *gorm.DB, videoid int64) (commentlist []Comment, nums int64) {
	db.Where("video_id = ? ", videoid).Find(&commentlist)
	nums = int64(len(commentlist))
	return
}

//使用userId查看点赞列表
func FavoriteQuerybyUserID(db *gorm.DB, userid int64) (videolist []Video, nums int64) {
	db.Model(&Video{}).Select("videos.*").Joins("inner join favorites on favorites.video_id = videos.id").Where("favorites.user_id = ? ", userid).Scan(&videolist)
	nums = int64(len(videolist))
	return
}

//查询用户(关注者)，作为user关注了哪些人的列表
func UserQueryByFollowID(db *gorm.DB, id int64) (users []User) {
	db.Model(&User{}).Select("users.*").Joins("inner join relations on relations.to_user_id = users.id").Where("relations.user_id = ? ", id).Scan(&users)
	return
}

//查询用户(关注者)，作为user被哪些人关注了
func UserQueryByFollowerID(db *gorm.DB, id int64) (users []User) {
	db.Model(&User{}).Select("users.*").Joins("inner join relations on relations.user_id = users.id").Where("relations.to_user_id = ? ", id).Scan(&users)
	return
}
