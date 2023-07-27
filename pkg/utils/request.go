package utils

import (
	"botserver/conf"
	"github.com/bytedance/sonic"
	"github.com/kaiheila/golang-bot/api/helper"
	log "github.com/sirupsen/logrus"
)

func SendMessage(
	msgType int,
	targetid string,
	content string,
	quote string,
	nonce string,
	tempIds string,
) error {
	//统一发送频道消息方法

	client := helper.NewApiHelper("/v3/message/create", conf.Token, conf.BaseUrl, "", "")
	data := map[string]interface{}{
		"type":           msgType,
		"target_id":      targetid,
		"content":        content,
		"quote":          quote,
		"nonce":          nonce,
		"temp_target_id": tempIds,
	}
	//将map转化成[]byte
	byteDtate, err := sonic.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := client.SetBody(byteDtate).Post()
	log.Info("正在发送频道消息:%s", client.String())
	if err != nil {
		log.Error("处理发送频道消息失败!", err)
	}
	log.Infof("发送频道成功,DATA:%s", string(resp))

	return nil
}

func UpdateMessage(id string, context string, quote string, tempIds string) error {
	//对频道信息进行更新
	//client := helper.NewApiHelper("/v3/message/update", conf.Token, conf.BaseUrl, "", "")

	return nil
}

func DeleteMessage(id string) error {
	//删除信息
	return nil
}

func GetChannel(id string) error {
	//获取用户所在语音频道id
	//get请求获取当前用户所在的语音频道，如果无则返回信息
	return nil
}
