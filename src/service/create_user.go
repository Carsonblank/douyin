package service

import (
	"errors"
	"github.com/RaymondCode/simple-demo/src/repository"
	"log"
)

type CreateUserFlow struct {
	userId   int64
	username string
	password string
}

func NewCreateUserFlow(username string, password string) *CreateUserFlow {
	return &CreateUserFlow{
		username: username,
		password: password,
	}
}

func CreateUser(username string, password string) (int64, error) {
	return NewCreateUserFlow(username, password).Do()
}

func (f *CreateUserFlow) Do() (int64, error) {
	if err := f.checkParam(); err != nil {
		return 0, err
	}
	if err := f.create(); err != nil {
		return 0, err
	}
	return f.userId, nil
}

func (f *CreateUserFlow) checkParam() error {
	//检查f.username和f.password是否合法,暂时不做
	user, _ := repository.NewUserDaoInstance().QueryUserByName(f.username)
	if user != nil {
		return errors.New("该用户已经存在")
	}
	return nil
}

func (f *CreateUserFlow) create() error {
	user := &repository.User{
		Name:     f.username,
		Password: f.password,
	}
	if err := repository.NewUserDaoInstance().CreateUser(user); err != nil {
		log.Println("insert err", err)
		return err
	}
	f.userId = user.Id
	return nil
}
