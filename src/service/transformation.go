package service

import (
	"douyin/src/database"
	"douyin/src/repository"
)

//根据提供的[]dtb.Video转化为可以用于message的video类型
//videoList_db：提供的dtb格式的video数据，userId：用户id，用来判断该视频对该用户来讲是不是已赞
func FromDBVideosToMesVideos(videoList_db []database.Video, userId int64) []Video {
	videos := make([]Video, 0, len(videoList_db))
	for _, v := range videoList_db {
		user_db, _ := repository.UserQueryByID(v.UserId)
		videos = append(videos, Video{
			Id:            v.Id,
			Author:        FromDBUsersTOMesUsers(user_db, userId)[0],
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.Favorite_count,
			CommentCount:  v.Comment_count,
			IsFavorite:    repository.FavoriteQueryByUserAndVideo(userId, v.Id),
			Title:         v.Title,
		})
	}
	return videos
}

//根据提供的[]dtb.User转化为可以用于message的user类型;用户id，用来判断用户userId是不是已经关注了列表中的用户
func FromDBUsersTOMesUsers(userList_db []database.User, userId int64) []User {
	users := make([]User, 0, len(userList_db))
	for _, u := range userList_db {
		users = append(users, User{
			Id:            u.Id,
			Name:          u.Name,
			FollowCount:   u.Follow_count,
			FollowerCount: u.Follower_count,
			IsFollow:      repository.RelationQueryByUserAndUser(userId, u.Id),
			Avatar:        u.Avatar,
			Signature:     u.Signature,
			TotalFavorite: u.Total_favorite,
			FavoriteCount: u.Favorite_count,
		})

	}
	return users
}
