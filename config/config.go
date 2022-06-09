package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/ini.v1"
)

var DSN string //= "douyin:123456@tcp(:3306)/douyindata?charset=utf8&parseTime=True&loc=Local" //数据库登陆DSN

var DouyinPort string //= "localhost:8080" //服务开放端口

var GIN_MODE string //gin的打开模式，release或debug

var MustVideosNums int //= 30 //返回视频流最大视频数量

type MySQLInfo struct {
	USER      string
	PASSWORD  string
	DB_HOST   string
	DB_PORT   string
	DB_NAME   string
	CHARSET   string
	ParseTime string
	Loc       string
}

type ServerInfo struct {
	HTTP_PORT string
	HTTP_HOST string
	MODE      string
}

//读取config.ini文件，构造端口号和打开模式
func GetServerConfig() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v \n", err)
		os.Exit(1)
	}
	d := new(ServerInfo)
	err = cfg.Section("server").MapTo(d)
	if err != nil {
		log.Printf("Fail to map server: %v \n", err)
		os.Exit(1)
	}
	DouyinPort = fmt.Sprintf("%s:%s", d.HTTP_HOST, d.HTTP_PORT)
	GIN_MODE = d.MODE
}

//构造DSN
func GetMySQLConfig() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v \n", err)
		os.Exit(1)
	}
	d := new(MySQLInfo)
	err = cfg.Section("mysql").MapTo(d)
	if err != nil {
		log.Printf("Fail to map MySQL: %v \n", err)
		os.Exit(1)
	}
	//DSN = "douyin:123456@tcp(127.0.0.1:3306)/douyindata?charset=utf8&parseTime=True&loc=Local"
	DSN = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%s&loc=%s",
		d.USER,
		d.PASSWORD,
		d.DB_HOST,
		d.DB_PORT,
		d.DB_NAME,
		d.CHARSET,
		d.ParseTime,
		d.Loc,
	)

}

type OtherInfo struct {
	MAX_FEED_VIDEO_NUMS int
}

//其他信息
func GetFeedNums() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v \n", err)
		os.Exit(1)
	}
	d := new(OtherInfo)
	err = cfg.Section("otherInfo").MapTo(d)
	if err != nil {
		log.Printf("Fail to map OtherInfo: %v \n", err)
		os.Exit(1)
	}
	MustVideosNums = d.MAX_FEED_VIDEO_NUMS
}
