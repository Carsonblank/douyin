package repository

import (
	"github.com/RaymondCode/simple-demo/src/database"
	"gorm.io/gorm"
	"log"
	"sync"
)

var db = database.MySqlDb

type User struct {
	Id            int64  `gorm:"column:id"`
	Name          string `gorm:"column:name"`
	Password      string `gorm:"column:password"`
	FollowCount   int64  `gorm:"column:follow_count"`
	FollowerCount int64  `gorm:"column:follower_count"`
}

func (User) TableName() string {
	return "user"
}

type UserDao struct {
}

var userDao *UserDao //空结构体节省内存
var userOnce sync.Once

func NewUserDaoInstance() *UserDao {
	//Do里面的函数只会执行一次
	userOnce.Do(
		func() {
			userDao = &UserDao{}
		})
	return userDao
}

func (*UserDao) CreateUser(user *User) error {
	if err := db.Create(user).Error; err != nil {
		log.Println("insert err", err.Error())
		return err
	}
	return nil
}

func (*UserDao) QueryUserByName(name string) (*User, error) {
	var user User
	err := db.Where("name = ?", name).Find(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, err
	}
	if err != nil {
		log.Println("find user by name err:", name, err)
		return nil, err
	}
	return &user, nil
}
