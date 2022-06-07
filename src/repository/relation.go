package repository

import (
	"douyin/src/database"
	db "douyin/src/database"

	"gorm.io/gorm"
)

//添加关注信息
func RelationCreate(rel db.Relation) (err error) {
	err = db.MySqlDb.Create(&rel).Error
	return
}

//删除关注信息
func RelationDelete(userid, touserid int64) (err error) {
	err = db.MySqlDb.Where("user_id = ? AND to_user_id = ?", userid, touserid).Delete(&db.Relation{}).Error
	return
}

//查询是否关注
func RelationQueryByUserAndUser(userid, touserid int64) (exist bool) {
	var nums int64
	if db.MySqlDb.Table("relations").Select("count(*)").Where("user_id = ? AND to_user_id = ?", userid, touserid).Scan(&nums); nums > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//查询用户(关注者)，作为user关注了哪些人的列表
func UserQueryByFollowID(id int64) (users []db.User) {
	db.MySqlDb.Model(&db.User{}).Select("users.*").Joins("inner join relations on relations.to_user_id = users.id").Where("relations.user_id = ? ", id).Scan(&users)
	return
}

//查询用户(关注者)，作为user被哪些人关注了
func UserQueryByFollowerID(id int64) (users []db.User) {
	db.MySqlDb.Model(&db.User{}).Select("users.*").Joins("inner join relations on relations.user_id = users.id").Where("relations.to_user_id = ? ", id).Scan(&users)
	return
}

//添加/取消关注操作，统一更新用户关注数、用户被关注数
func RelationUpdataNumbers(userId, toUserId int64, add bool) error {
	var n int64
	if add {
		n = 1
	} else {
		n = -1
	}
	//更新用户关注数和被关注数
	if err := database.MySqlDb.Model(&database.User{}).Where("id = ?", userId).Update("follow_count", gorm.Expr("follow_count + ?", n)).Error; err != nil {
		return err
	}
	if err := database.MySqlDb.Model(&database.User{}).Where("id = ?", toUserId).Update("follower_count", gorm.Expr("follower_count + ?", n)).Error; err != nil {
		return err
	}
	return nil
}
