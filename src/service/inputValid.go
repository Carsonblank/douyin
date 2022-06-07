package service

import (
	"strconv"
	"time"
)

//定义的error
type MyError struct {
	ErrMsg string
}

func (s MyError) Error() string {
	return s.ErrMsg
}

//判断token是否合法，对于不合法的token返回的Id=0
func Token2Id(tokenString string) (userId int64) {
	if id, isValid := TokenValidity(tokenString); !isValid {
		userId = 0
	} else {
		//合法的token解析出登陆用户ID并赋值给token_user_id变量
		userId = id
	}
	return
}

//验证feed接口给定的时间戳是否合法，对于不合法的时间戳和不给定时间戳的情况，时间戳为当前时间
func TimeStamp(latestTimeString string) (lastTime int64) {
	lastTime, _ = strconv.ParseInt(latestTimeString, 10, 64)
	//排除非法的lastTime输入
	if lastTime <= 0 {
		lastTime = time.Now().UnixMilli()
	}
	return
}

//验证赞操作的操作数是否合法
func FavoriteActionNum(actionType string) (int64, error) {
	at, _ := strconv.ParseInt(actionType, 10, 64)
	if at != 1 && at != 2 {
		return 0, MyError{"Favorite action type error."}
	}
	return at, nil
}

//验证评论操作的操作数是否合法
func CommentActionNum(actionType string) (int64, error) {
	at, _ := strconv.ParseInt(actionType, 10, 64)
	if at != 1 && at != 2 {
		return 0, MyError{"Comment action type error."}
	}
	return at, nil
}

//验证关注操作的操作数是否合法
func RelationActionNum(actionType string) (int64, error) {
	at, _ := strconv.ParseInt(actionType, 10, 64)
	if at != 1 && at != 2 {
		return 0, MyError{"Relation action type error."}
	}
	return at, nil
}
