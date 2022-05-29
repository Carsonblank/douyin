//提供一些基础的函数
package controller

import (
	dtb "douyin/database"

	"gorm.io/gorm"
)

//根据提供的[]dtb.Video转化为可以用于message的video类型
//videoList_db：提供的dtb格式的video数据，userId：用户id，用来判断该视频对该用户来讲是不是已赞
func FromDBVideosToMesVideos(db *gorm.DB, videoList_db []dtb.Video, userId int64) []Video {
	videos := make([]Video, 0, len(videoList_db))
	for _, v := range videoList_db {
		user_db, _ := dtb.UserQueryByID(db, v.UserId)
		videos = append(videos, Video{
			Id:            v.Id,
			Author:        FromDBUsersTOMesUsers(db, user_db, userId)[0],
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: dtb.VideoFavoriteCount(db, v.Id),
			CommentCount:  dtb.VideoCommentCount(db, v.Id),
			IsFavorite:    dtb.FavoriteQueryByUserAndVideo(db, userId, v.Id),
			Title:         v.Title,
		})
	}
	return videos
}

//根据提供的[]dtb.User转化为可以用于message的user类型;用户id，用来判断用户userId是不是已经关注了列表中的用户
func FromDBUsersTOMesUsers(db *gorm.DB, userList_db []dtb.User, userId int64) []User {
	users := make([]User, 0, len(userList_db))
	for _, u := range userList_db {
		users = append(users, User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   dtb.UserFollowCount(db, u.Id),
			FollowerCount: dtb.UserFollowerCount(db, u.Id),
			IsFollow:      dtb.FavoriteQueryByUserAndUser(db, userId, u.Id),
			Avatar:        u.Avatar,
			Signature:     u.Signature,
			TotalFavorite: dtb.UserFavoritedCount(db, u.Id),
			FavoriteCount: dtb.UserFavoriteCount(db, u.Id),
		})

	}
	return users
}

//根据用户名和密码生成token
func GetToken(username, password string) (token string) {
	token = username + password
	return
}
