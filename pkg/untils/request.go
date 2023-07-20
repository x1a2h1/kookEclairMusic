package untils

import (
	"botserver/conf"
	"github.com/bytedance/sonic"
	"github.com/idodo/golang-bot/kaihela/api/helper"
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
