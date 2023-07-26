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
	"sync"
)

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

func GetChannelId(gid string, uid string) (string, error) {

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
	return cid, err
}
func NewClient(token string, channelId string) (*VoiceInstance, error) {
	vi := VoiceInstance{
		Token:     token,
		ChannelId: channelId,
	}
	err := vi.Init()
	if err != nil {
		return nil, err
	}
	return &vi, nil
}

var Status sync.Map

func Play(gid string, cid string, uid string) error {
	// 通过cid channel ID创建播放
	//判断当前传进来的id是否活跃状态
	//判断当前频道id是否已经创建了持久化进程
	isPlay, ok := Status.Load(cid)
	if !ok {
		fmt.Println(isPlay)
		Status.Store(cid, true)
		client, err := NewClient(conf.Token, cid)
		if err != nil {
			return err
		}
		var playlist model.Playlist
		conf.DB.Preload("Songs").Find(&playlist, gid)
		go func() {
			client.Init()
			defer client.Close()
			defer client.wsConnect.Close()
			for {
				songInfo := getMusic(gid)
				if songInfo.ID == 0 {
					break
				}
				fmt.Println("当前正在播放歌曲", songInfo.SongID)
				url, times := song.GetMusicUrl(songInfo.SongID)
				err := client.PlayMusic(url)
				if err != nil {
					return
				}
				conf.DB.Debug().Delete(&songInfo, songInfo.ID)
				fmt.Println("当前歌曲："+songInfo.Name+"，总用时：", times)

			}
			fmt.Println(cid, "频道播放已结束")
			Status.Delete(cid)
			fmt.Println(cid, "频道播放已结束！进程退出成功！")
		}()
		//	goroutine结束后
		fmt.Println("已经开启goroutine进行连接播放")
	} else {
		fmt.Println("当前频道" + cid + "播放列表正在播放")
	}

	return nil
}

type Song struct {
	ID       int
	SongID   string
	Name     string
	CoverUrl string
	UserName string
	Singer   string
}

func getMusic(gid string) Song {
	var Playlist model.Playlist
	conf.DB.Preload("Songs").Find(&Playlist, gid)
	if len(Playlist.Songs) == 0 {
		return Song{
			ID:       0,
			SongID:   "",
			Name:     "",
			CoverUrl: "",
			UserName: "",
			Singer:   "",
		}
	} else {
		SongInfo := Song{
			ID:       Playlist.Songs[0].ID,
			SongID:   Playlist.Songs[0].SongId,
			Name:     Playlist.Songs[0].SongName,
			CoverUrl: Playlist.Songs[0].CoverUrl,
			UserName: Playlist.Songs[0].UserName,
			Singer:   Playlist.Songs[0].UserName,
		}
		return SongInfo
	}
}
