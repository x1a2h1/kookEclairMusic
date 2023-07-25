package kook

import (
	model "botserver/app/model"
	"botserver/app/song"
	"botserver/conf"
	"botserver/pkg/untils"
	"encoding/json"
	"fmt"
	"github.com/kaiheila/golang-bot/api/helper"
	log "github.com/sirupsen/logrus"
)

var MusicList = make(chan model.MusicList, 50)

func init() {

}

func PlayForList(gid string, targerId string) error {
	//	根据播放列表播放
	var Current = make(map[string]string)
	Current["songName"] = "陈宇晖牛逼"
	Current["songUrl"] = "https://c-ssl.dtstatic.com/uploads/blog/202207/09/20220709150824_97667.thumb.400_0.jpg"
	Current["singer"] = "陈宇晖"
	Current["userName"] = "夏至"
	var playlist model.Playlist
	conf.DB.Preload("Songs").Take(&playlist, gid)

	listMsg := model.CardMessageCard{
		Theme: model.CardThemeSuccess,
		Modules: []interface{}{
			&model.CardMessageHeader{Text: model.CardMessageElementText{
				Content: "播放列表",
				Emoji:   false,
			}},
			&model.CardMessageDivider{},
			&model.CardMessageSection{
				Mode: "right",
				Text: model.CardMessageElementKMarkdown{Content: "> ** **\n> **" + Current["songName"] + "**\n> **歌手**\n> ** **"},
				Accessory: &model.CardMessageElementImage{
					Src:    Current["songUrl"],
					Size:   "sm",
					Circle: true,
				},
			},
			&model.CardMessageDivider{},
		},
	}
	for _, item := range playlist.Songs {
		listMsg.AddModule(
			&model.CardMessageSection{
				Mode: model.CardMessageSectionModeLeft,
				Text: model.CardMessageElementKMarkdown{
					Content: "> ** **\n> **" + item.SongName + "**\n> **歌手**\n> ** **",
				},
				Accessory: &model.CardMessageElementImage{
					Src:  item.CoverUrl,
					Size: "sm",
				},
			})
	}
	sendMsg, err := model.CardMessage{&listMsg}.BuildMessage()
	if err != nil {
		return err
	}

	err = untils.SendMessage(10, targerId, sendMsg, "", "", "")
	if err != nil {
		return err
	}

	fmt.Println("当前所在的服务器列表", playlist)
	return err
}

type GateWayHttpApiResult struct {
	Code    int32  `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Url string `json:"url"`
	} `json:"data"`
}

func SongsHandle(songs *model.MusicList) {
	fmt.Println("403335371往Channel发送消息")
	songUrl, times := song.GetMusicUrl(songs.SongId)
	sendData := model.MusicList{
		Guild:    songs.Guild,
		ChanId:   songs.ChanId,
		SongId:   songs.SongId,
		SongName: songs.SongName,
		MusicUrl: songUrl,
		UserName: songs.UserName,
		CoverUrl: songs.CoverUrl,
		Duration: times,
	}
	MusicList <- sendData
}
func MusicHandle() {
	//	接收当前服务器，当前频道id，以及用户名进行处理
	//	对数据库进行相对应的检索
	//接收当前服务器是否存在播放，如果存在播放for if判断
	//	将当前歌曲传入channel
}

func GetChannelId(gid string, uid string) string {

	Client := helper.NewApiHelper("/v3/channel-user/get-joined-channel", conf.Token, conf.BaseUrl, "", "")
	Client.SetQuery(map[string]string{
		"guild_id": gid,
		"user_id":  uid,
	})
	resp, err := Client.Get()
	log.Info("正在发送频道消息:%s", Client.String())
	if err != nil {
		log.Error("处理发送频道消息失败!", err)
	}
	log.Infof("发送频道成功,DATA:%s", string(resp))
	type Response struct {
		Data struct {
			Items []struct {
				Id string
			}
		}
	}
	var res Response
	err = json.Unmarshal(resp, &res)
	cid := res.Data.Items[0].Id
	return cid
}
