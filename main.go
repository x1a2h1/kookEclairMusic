package main

import (
	"botserver/app"
	"botserver/conf"
	"fmt"
	"github.com/kaiheila/golang-bot/api/base"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
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
	sendBotOnline(conf.OnlineUUID)
	session.Start()
	ticker := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-ticker.C:
			sendBotOnline(conf.OnlineUUID)
		}
	}
}

func sendBotOnline(uuid string) {
	req, err := http.NewRequest("GET", "http://bot.gekj.net/api/v1/online.bot", nil)
	if err != nil {
		log.Error(err)
	}
	req.Header.Set("uuid", uuid)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error("发送请求错误:", err)
		return
	}
	defer resp.Body.Close()

	// 这里可以处理响应，例如打印 HTTP 状态码
	log.Println("Response status:", resp.Status)
}
