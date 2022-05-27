package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/ini.v1"
	"log"
	"os"
	"time"
)

type MySQLInfo struct {
	TYPE         string
	USER         string
	PASSWORD     string
	DB_HOST      string
	DB_PORT      string
	DB_NAME      string
	CHARSET      string
	ParseTime    string
	MaxIdleConns int
	MaxOpenConns int
}

type ServerInfo struct {
	HTTP_PORT string
	HTTP_HOST string
	MODE      string
}

func GetServerConfig() *ServerInfo {
	cfg, err := ini.Load("configFile/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v \n", err)
		os.Exit(1)
	}
	d := new(ServerInfo)
	_ = cfg.Section("server").MapTo(d)
	return d
}

func GetMySQLConfig() *MySQLInfo {
	cfg, err := ini.Load("configFile/config.ini")
	if err != nil {
		log.Printf("Fail to read file: %v \n", err)
		os.Exit(1)
	}
	d := new(MySQLInfo)
	_ = cfg.Section("mysql").MapTo(d)
	return d
}

func GetLogPath() string {
	timeObj := time.Now()
	datetime := timeObj.Format("2006-01-02-15-04-05")
	return "logfile/Course_Select" + datetime + ".log"
}

func GetLogFormat(param gin.LogFormatterParams) string {
	return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.Request.Proto,
		param.StatusCode,
		param.Latency,
		param.Request.UserAgent(),
		param.ErrorMessage,
	)
}
