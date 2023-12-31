package main

import (
	"botserver/app"
	"botserver/conf"
	"botserver/pkg/utils"
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
	ticker := time.NewTicker(30 * time.Minute)
	session := base.NewWebSocketSession(conf.Token, conf.BaseUrl, "./session.pid", "", 1)
	session.On(base.EventReceiveFrame, &app.ReceiveFrameHandler{})
	session.On("GROUP*", &app.GroupEventHandler{})
	session.On("GROUP_9", &app.GroupTextEventHandler{Token: conf.Token, BaseUrl: conf.BaseUrl})
	sendBotOnline(conf.OnlineUUID)
	go func() {
		for {
			select {
			case <-ticker.C:
				sendBotOnline(conf.OnlineUUID)
			}
		}
	}()
	session.Start()

}

func sendBotOnline(uuid string) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://bot.gekj.net/api/v1/online.bot", nil)
	if err != nil {
		log.Error("机器人在线验证出现未知问题", err)
		return
	}
	req.Header.Set("uuid", uuid)
	resp, err := client.Do(req)
	if err != nil {
		log.Error("发送请求错误:", err)
		return
	}
	defer resp.Body.Close()
	// 这里可以处理响应，例如打印 HTTP 状态码
	log.Println("请求状态:", resp.Body)
	utils.SendMessage(1, "1518342897030347", "在线验证成功！", "", "", "")

}
