package controller

import (
	"douyin/common"
	dtb "douyin/database"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//点赞操作
func FavoriteAction(c *gin.Context) {

	token := c.Query("token")
	var token_user_id int64
	//查看token是否合法
	if id, isValid := common.TokenValidity(token); !isValid {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "Token is invalid.",
		},
		)
		//对于点赞操作，token非法直接返回
		return
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		token_user_id = id
		//查看合法token返回的用户ID是否存在
		if _, exist := dtb.UserQueryByID(common.DB, token_user_id); !exist {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "Login User doesn't exist",
			},
			)
			return
		}
	}

	video_id, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	action_type, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)
	var video_user_id int64
	//查看该视频是否存在
	if exist, id := dtb.VideoQueryByID(common.DB, video_id); !exist {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1, StatusMsg: "Video is not exist!\n",
		})
		return
	} else {
		video_user_id = id
	}
	//添加点赞记录(action_type==1)和去除点赞记录(action_type==2)
	if action_type == 1 {
		if exist := dtb.FavoriteQueryByUserAndVideo(common.DB, token_user_id, video_id); exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat favorite action."})
		} else {
			if err := dtb.FavoriteCreate(common.DB, dtb.Favorite{
				UserId:   token_user_id,
				ToUserId: dtb.UserQueryByVideoID(common.DB, video_id),
				VideoId:  video_id,
			}); err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
				})
			} else {
				//更新视频点赞数
				if err := common.DB.Model(&dtb.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err == nil {
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				} else {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
					})
				}
				//更新video用户获赞数
				if err := common.DB.Model(&dtb.User{}).Where("id = ?", video_user_id).Update("total_favorite", gorm.Expr("total_favorite + ?", 1)).Error; err == nil {
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				} else {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
					})
				}
				//更新token用户喜欢数
				if err := common.DB.Model(&dtb.User{}).Where("id = ?", video_user_id).Update("favorite_count", gorm.Expr("favorite_count + ?", 1)).Error; err == nil {
					c.JSON(http.StatusOK, Response{StatusCode: 0})
				} else {
					c.JSON(http.StatusOK, Response{
						StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
					})
				}
			}
		}
	} else if action_type == 2 {
		//查看要取消的赞是否存在
		if exist := dtb.FavoriteQueryByUserAndVideo(common.DB, token_user_id, video_id); !exist {
			c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Repeat delete favorite action."})
			return
		}
		err := dtb.FavoriteDelete(common.DB, token_user_id, video_id)
		if err != nil { //取消赞失败
			c.JSON(http.StatusOK, Response{
				StatusCode: 1, StatusMsg: fmt.Sprintf("Remove favorite action error : %v\n", err),
			})
		} else { //取消赞成功
			//更新点赞数
			if err := common.DB.Model(&dtb.Video{}).Where("id = ?", video_id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err == nil {
				c.JSON(http.StatusOK, Response{StatusCode: 0})
			} else {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
				})
			}
			//更新用户获赞数
			if err := common.DB.Model(&dtb.User{}).Where("id = ?", video_user_id).Update("total_favorite", gorm.Expr("total_favorite - ?", 1)).Error; err == nil {
				c.JSON(http.StatusOK, Response{StatusCode: 0})
			} else {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
				})
			}
			//更新用户喜欢数
			if err := common.DB.Model(&dtb.User{}).Where("id = ?", token_user_id).Update("favorite_count", gorm.Expr("favorite_count - ?", 1)).Error; err == nil {
				c.JSON(http.StatusOK, Response{StatusCode: 0})
			} else {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1, StatusMsg: fmt.Sprintf("Favorite action error : %v\n", err),
				})
			}
		}
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "Favorite action_type must be 1 or 2"})
	}

}

//获取点赞列表
func FavoriteList(c *gin.Context) {
	user_id, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	if user, exist := dtb.UserQueryByID(common.DB, user_id); exist {
		videoList, _ := dtb.FavoriteQuerybyUserID(common.DB, user_id)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: FromDBVideosToMesVideos(common.DB, videoList, user[0].Id),
		})
	} else {
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 1,
			},
		})
	}

}
