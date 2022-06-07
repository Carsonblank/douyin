package repository

import (
	"douyin/src/database"
	db "douyin/src/database"

	"gorm.io/gorm"
)

//添加评论信息
func CommentCreate(com *db.Comment) (err error) {
	err = db.MySqlDb.Create(com).Error
	return
}

//使用videoId查询评论列表
func CommentQuerybyVideoID(videoid int64) (commentlist []db.Comment, nums int64) {
	db.MySqlDb.Where("video_id = ? ", videoid).Find(&commentlist)
	nums = int64(len(commentlist))
	return
}

//查询评论是否存在
func CommentQueryByCommentId(comId int64) bool {
	var nums int64
	if db.MySqlDb.Table("comments").Select("count(*)").Where("id = ?", comId).Scan(&nums); nums > 0 {
		return true
	}
	return false

}

//删除评论
func CommentDelete(comment_id int64) (err error) {
	err = db.MySqlDb.Where("id = ? ", comment_id).Delete(&db.Comment{}).Error
	return
}

//添加/取消评论操作
//添加/取消赞操作，统一更新视频点赞数、用户获赞数、用户喜欢数
func CommentUpdataNumbers(videoId int64, add bool) error {
	var n int64
	if add {
		n = 1
	} else {
		n = -1
	}
	//更新视频评论数
	if err := database.MySqlDb.Model(&database.Video{}).Where("id = ?", videoId).Update("comment_count", gorm.Expr("comment_count + ?", n)).Error; err != nil {
		return err
	}
	return nil
}
