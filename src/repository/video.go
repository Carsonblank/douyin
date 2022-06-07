package repository

import (
	"douyin/config"
	db "douyin/src/database"
)

//添加视频信息
func VideoCreate(video db.Video) (err error) {
	err = db.MySqlDb.Create(&video).Error
	return
}

//根据videoId查看video是否存在
func VideoQueryByID(video_id int64) (exist bool, user_id int64) {
	var v []db.Video
	db.MySqlDb.Where("id = ? ", video_id).Find(&v)
	if len(v) > 0 {
		exist = true
		user_id = v[0].UserId
	} else {
		exist = false
		user_id = 0
	}
	return

}

//返回视频上传时间不大于给定时间的最多nums个视频的列表，nums在config文件夹给定
func VideoQueryByLastTime(lastest_time int64) []db.Video {
	videolist_db := make([]db.Video, 0, config.MustVideosNums)
	//降序排列
	db.MySqlDb.Limit(config.MustVideosNums).Order("created_at DESC").Where("created_at <= ?", lastest_time).Find(&videolist_db)

	return videolist_db
}

//使用UserID查询视频列表
func VideoQueryByUserID(userID int64) (videolist []db.Video) {
	db.MySqlDb.Where("user_id = ? ", userID).Find(&videolist)
	return
}
