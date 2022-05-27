package controller

import (
	"github.com/RaymondCode/simple-demo/src/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

// Feed same demo video list for every request
func Feed(c *gin.Context) {
	videosFromRepository, _ := repository.NewVideoDaoInstance().ReturnAllVideos()
	n := len(videosFromRepository)
	returnVideos := make([]Video, n)
	for i := 0; i < n; i++ {
		returnVideos[i] = Video{
			Id:            videosFromRepository[i].Id,
			Author:        User{},
			PlayUrl:       videosFromRepository[i].PlayURL,
			CoverUrl:      videosFromRepository[i].CoverURL,
			FavoriteCount: videosFromRepository[i].FavoriteCount,
			CommentCount:  videosFromRepository[i].CommentCount,
			IsFavorite:    false,
		}
	}
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: returnVideos,
		NextTime:  time.Now().Unix(),
	})
}
