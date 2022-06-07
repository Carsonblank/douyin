package repository

import (
	db "douyin/src/database"
)

//向Users数据库添加用户信息
func UserCreate(user *db.User) (err error) {
	err = db.MySqlDb.Create(user).Error
	return
}

//使用主键id查询用户
func UserQueryByID(id int64) (user []db.User, exist bool) {
	if db.MySqlDb.Limit(1).Find(&user, id); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//使用Name查询用户
func UserQueryByName(username string) (user []db.User, exist bool) {
	if db.MySqlDb.Limit(1).Where("name = ?", username).Find(&user); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}

//验证登录信息是否正确(账户密码)
func UserValid(username, password string) (user []db.User, exist bool) {
	if db.MySqlDb.Limit(1).Where("name = ? AND password = ?", username, password).Find(&user); len(user) > 0 {
		exist = true
	} else {
		exist = false
	}
	return
}
