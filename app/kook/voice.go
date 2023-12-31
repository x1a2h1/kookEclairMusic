package kook

import (
	model "botserver/app/model"
	redisDB "botserver/app/redis"
	"botserver/app/song"
	"botserver/conf"
	"botserver/pkg/utils"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/kaiheila/golang-bot/api/helper"
	"github.com/nareix/joy4/format/rtmp"
	log "github.com/sirupsen/logrus"
	"github.com/x1a2h1/kookvoice"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
)

func init() {
}

// 当前正在播放的音乐
var CurrentSong sync.Map

type SongData struct {
	GuildID  string
	Singer   string
	SongName string
	UserName string
	CoverUrl string
}

func PlayForList(gid string, targerId string) error {
	//	根据播放列表播放
	var playlist model.Playlist
	conf.DB.Preload("Songs").Take(&playlist, gid)
	//获取当前频道正在播放的音乐
	SongInfoData := &SongData{}
	CurrentData, ok := CurrentSong.Load(gid)
	if !ok {
		SongInfoData = &SongData{
			GuildID:  gid,
			Singer:   " ",
			SongName: "当前暂无播放",
			UserName: " ",
			CoverUrl: "https://c-ssl.dtstatic.com/uploads/blog/202207/09/20220709150824_97667.thumb.400_0.jpg",
		}
	} else {
		if data, ok := CurrentData.(SongData); ok {
			SongInfoData = &SongData{
				GuildID:  data.GuildID,
				Singer:   data.Singer,
				SongName: data.SongName,
				UserName: data.UserName,
				CoverUrl: data.CoverUrl,
			}
		} else {
			err := utils.SendMessage(1, targerId, "当前列表存在未知错误！", "", "", "")
			if err != nil {
				return err
			}
		}
	}

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
				Text: model.CardMessageElementKMarkdown{Content: "> ** **\n> **" + SongInfoData.SongName + "**\n> **" + SongInfoData.Singer + "**\n "},
				Accessory: &model.CardMessageElementImage{
					Src:    SongInfoData.CoverUrl,
					Size:   "sm",
					Circle: true,
				},
			},
			&model.CardMessageDivider{},
		},
	}
	for index, item := range playlist.Songs {
		if index > 20 {
			break
		}
		listMsg.AddModule(
			&model.CardMessageSection{
				Mode: model.CardMessageSectionModeLeft,
				Text: model.CardMessageElementKMarkdown{
					Content: "> ** **\n> **" + item.SongName + "**\n> **" + item.SongSinger + "**\n> ** **",
				},
				Accessory: &model.CardMessageElementImage{
					Src:  item.CoverUrl,
					Size: "sm",
				},
			})
	}
	listMsg.AddModule(
		&model.CardMessageSection{Text: model.CardMessageElementText{Content: "当前卡片仅展示20首歌曲"}},
	)
	sendMsg, err := model.CardMessage{&listMsg}.BuildMessage()
	if err != nil {
		return err
	}

	err = utils.SendMessage(10, targerId, sendMsg, "", "", "")
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
	if len(res.Data.Items) == 0 {
		return "", err
	} else {
		cid := res.Data.Items[0].Id
		return cid, err
	}

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
var Cmds = make(map[string]*exec.Cmd)
var Mu sync.Mutex
var TotalPlay int32

func Play(gid string, cid string, uid string) error {
	// 通过cid channel ID创建播放
	//判断当前传进来的id是否活跃状态
	//判断当前频道id是否已经创建了持久化进程
	isPlay, ok := Status.Load(gid)
	if !ok {
		fmt.Println(isPlay)
		Status.Store(gid, true)
		var playlist model.Playlist
		conf.DB.Preload("Songs").Find(&playlist, gid)
		go func(GuildID string) {
			//开始播放当前服务器列表歌曲
			for {

				songInfo := getMusic(GuildID)
				if songInfo.ID == 0 {
					fmt.Println("当前服务器中已没有可播放歌曲\n")
					break
				}
				//TODO 获取当前语音频道id中的数据，如果没人就清空列表并break
				//TODO 获取当前语音频道id中的数据，如果没人就break结束
				//TODO 请求当前服务器的cid本地推流地址
				//TODO 接收本地推流发送至kook官方地址
				gatewayUrl := kookvoice.GetGatewayUrl(conf.Token, cid)
				connect, rtpUrl := kookvoice.InitWebsocketClient(gatewayUrl)
				defer connect.Close()
				go kookvoice.KeepWebsocketClientAlive(connect)
				go kookvoice.KeepRecieveMessage(connect)
				CurrentSong.Store(gid, SongData{
					GuildID:  GuildID,
					Singer:   songInfo.Singer,
					SongName: songInfo.Name,
					UserName: songInfo.UserName,
					CoverUrl: songInfo.CoverUrl,
				})
				url, times := song.GetMusicUrl(songInfo.SongID)
				fmt.Println("当前播放的音乐url为", url, "当前播放时长为", times)
				conf.DB.Debug().Delete(&songInfo, songInfo.ID)
				if url == "" || times == 0 {
					log.Error("获取音乐url失败")
					continue
				}
				//判断播放结束、时长为0，链接不存在关闭ffmpeg进程
				fmt.Println("\n\n开始播放歌曲", songInfo.Name, "歌曲id", songInfo.SongID)
				//kookvoice.StreamAudio(rtpUrl, url)

				StreamAudio(gid, rtpUrl, url)
				//播放结束，已播放数量加1
				err := redisDB.Rdb.Incr(context.Background(), "totalPlayed").Err()
				if err != nil {
					continue
				}
				atomic.AddInt32(&TotalPlay, 1)
				fmt.Println("歌曲："+songInfo.Name+"，总用时：", times, "\n\n>>>播放结束<<<")
			}
			CurrentSong.Delete(gid)
			Status.Delete(gid)
			fmt.Println(cid, "频道播放已结束！进程退出成功！")
		}(gid)
		//	goroutine结束后
		fmt.Println("已经开启goroutine进行连接播放")
	} else {
		fmt.Println("当前服务器" + gid + "播放列表正在播放，已将歌曲添加至列表")
	}
	return nil
}

type voice struct {
	Token         string
	ChannelId     string
	wsConnect     *websocket.Conn
	streamProcess *os.Process
	sourceProcess *os.Process
}

func (i *voice) Init() string {
	gatewayUrl := kookvoice.GetGatewayUrl(i.Token, i.ChannelId)
	connect, rtpUrl := kookvoice.InitWebsocketClient(gatewayUrl)

	go kookvoice.KeepWebsocketClientAlive(connect)
	go kookvoice.KeepRecieveMessage(connect)
	i.wsConnect = connect
	return rtpUrl
}

func localStream(gid string) {
	//TODO 创建本地推流地址
	server := &rtmp.Server{Addr: ":6324"}

	server.HandlePublish = func(conn *rtmp.Conn) {
		for {
			songInfo := getMusic(gid)
			if songInfo.ID == 0 {
				break
			}
		}
	}
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
			Singer:   Playlist.Songs[0].SongSinger,
		}
		return SongInfo
	}
}
