package main

import (
	"botserver/app"
	"botserver/conf"
	"fmt"
	"github.com/idodo/golang-bot/kaihela/api/base"
	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetReportCaller(true)
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	fmt.Println("原神启动！")
	session := base.NewWebSocketSession(conf.Token, conf.BaseUrl, "./session.pid", "", 1)
	session.On(base.EventReceiveFrame, &app.ReceiveFrameHandler{})
	session.On("GROUP*", &app.GroupEventHandler{})
	session.On("GROUP_9", &app.GroupTextEventHandler{Token: conf.Token, BaseUrl: conf.BaseUrl})
	session.Start()
}