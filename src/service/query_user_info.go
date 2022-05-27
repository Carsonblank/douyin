package service

import (
	"github.com/RaymondCode/simple-demo/src/repository"
	"log"
)

type UserInfo struct {
	User *repository.User
}

type QueryUserInfoFlow struct {
	username string
	userInfo *UserInfo
}

func QueryUserInfo(name string) (*UserInfo, error) {
	return NewQueryUserInfoFlow(name).Do()
}

func NewQueryUserInfoFlow(name string) *QueryUserInfoFlow {
	return &QueryUserInfoFlow{
		username: name,
	}
}
func (f *QueryUserInfoFlow) Do() (*UserInfo, error) {
	/*
		if err := f.checkParam(); err != nil {
			return nil, err
		}
		if err := f.prepareInfo(); err != nil {
			return nil, err
		}
		if err := f.packPageInfo(); err != nil {
			return nil, err
		}

	*/
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	return f.userInfo, nil
}
func (f *QueryUserInfoFlow) prepareInfo() error {
	user, err := repository.NewUserDaoInstance().QueryUserByName(f.username)
	if err != nil {
		log.Println(err)
		return err
	}
	f.userInfo = &UserInfo{User: user}
	return nil
}
