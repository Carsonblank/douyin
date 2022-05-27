package repository

import (
	"log"
	"sync"
)

type Video struct {
	Id            int64  `gorm:"column:id"`
	UserID        string `gorm:"column:userID"`
	PlayURL       string `gorm:"column:play_url"`
	CoverURL      string `gorm:"column:cover_url"`
	FavoriteCount int64  `gorm:"column:favorite_count"`
	CommentCount  int64  `gorm:"column:comment_count"`
}

func (Video) TableName() string {
	return "video"
}

type VideoDao struct {
}

var videoDao *VideoDao //空结构体节省内存
var videoOnce sync.Once

func NewVideoDaoInstance() *VideoDao {
	//Do里面的函数只会执行一次
	videoOnce.Do(
		func() {
			videoDao = &VideoDao{}
		})
	return videoDao
}

func (*VideoDao) CreateVideo(video *Video) error {
	if err := db.Create(video).Error; err != nil {
		log.Println("CreateVideo err", err.Error())
		return err
	}
	return nil
}

func (*VideoDao) ReturnAllVideos() ([]*Video, error) {
	var videos []*Video
	err := db.Find(&videos).Error
	if err != nil {
		log.Println("ReturnAllVideos err:", err.Error())
		return nil, err
	}
	return videos, nil
}
