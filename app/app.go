package app

import (
	"botserver/app/kook"
	"botserver/app/model"
	_ "botserver/app/redis"
	"botserver/app/song"
	"botserver/conf"
	"botserver/pkg/utils"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	"github.com/kaiheila/golang-bot/api/base"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	log "github.com/sirupsen/logrus"
	"github.com/x1a2h1/kookvoice"
	"os"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

func init() {

}

// 定义一个统一发送消息的卡片

// 监听事件
type ReceiveFrameHandler struct {
}

func (rf *ReceiveFrameHandler) Handle(e event.Event) error {
	log.WithField("event", e).Info("PPAP监听中···")
	return nil
}

type GroupEventHandler struct {
}

func (ge *GroupEventHandler) Handle(e event.Event) error {
	log.WithField("event", e).Info("监听群事件···")
	return nil
}

type GroupTextEventHandler struct {
	Token   string
	BaseUrl string
}

func (gte *GroupTextEventHandler) Handle(e event.Event) error {
	log.WithField("event", fmt.Sprintf("%+v", e.Data())).Info("收到频道内的文字消息.")
	err := func() error {
		if _, ok := e.Data()[base.EventDataFrameKey]; !ok {
			return errors.New("数据丢失！")
		}
		frame := e.Data()[base.EventDataFrameKey].(*event2.FrameMap)
		data, err := sonic.Marshal(frame.Data)
		if err != nil {
			return err
		}
		msgEvent := &event2.MessageKMarkdownEvent{}
		err = sonic.Unmarshal(data, msgEvent)
		log.Infof("收到json事件:%+v", msgEvent)
		if err != nil {
			return err
		}
		if msgEvent.Author.Bot {
			log.Info("机器人消息")
			return nil
		}

		//开启多线程
		//开启多线程结束https://kook.top/x2eAZA
		LinkUrl := "https://kook.top/x2eAZA"
		helpCard := model.CardMessageCard{
			Theme: "none",
			Size:  model.CardSizeLg,
			Modules: []interface{}{
				&model.CardMessageHeader{Text: model.CardMessageElementText{
					Content: "🌈帮助菜单",
					Emoji:   true,
				}},
				&model.CardMessageDivider{},

				&model.CardMessageSection{
					Text: model.CardMessageParagraph{
						Cols: 3,
						Fields: []interface{}{
							model.CardMessageElementKMarkdown{Content: "**指令**\n(font)/网易 { 歌曲名 } (font)[error]\n/帮助\n/状态"},
							model.CardMessageElementKMarkdown{Content: "**功能**\n(font)播放网易云音乐(font)[success]\n帮助菜单\n当前机器人状态"},
							model.CardMessageElementKMarkdown{Content: "**示例**\n/网易 乐鼓 (dj版)\n/帮助\n/状态"},
						},
					},
				},
				&model.CardMessageInvite{
					Code: LinkUrl,
				},
				&model.CardMessageDivider{},
				&model.CardMessageContext{
					&model.CardMessageElementText{Content: "当前频道id：" + msgEvent.TargetId + "\n"},
					&model.CardMessageElementText{Content: "当前频道名：" + msgEvent.ChannelName + "\n"},
					&model.CardMessageElementText{Content: "当前频道服务器ID：" + msgEvent.GuildID + "\n"},
				},
				&model.CardMessageSection{
					Text: model.CardMessageElementKMarkdown{Content: "Version:" + "`" + conf.Version + "` 问题反馈(met)1260041158(met)"},
				},
			},
		}
		helpCardMsg, err := model.CardMessage{&helpCard}.BuildMessage()
		if err != nil {
			log.Error("编译信息时出错！", err)
		}
		fmt.Println(helpCardMsg)

		//InviteCard := model.CardMessage{
		//	&model.CardMessageCard{
		//		Theme: "",
		//		Modules: []interface{}{
		//
		//			&model.CardMessageHeader{Text: model.CardMessageElementText{
		//				Content: LinkUrl,
		//				Emoji:   false,
		//			}},
		//			&model.CardMessageInvite{
		//				Code: LinkUrl,
		//			},
		//		},
		//	},
		//}.MustBuildMessage()

		if msgEvent.Content == "/帮助" {
			go utils.SendMessage(10, msgEvent.TargetId, helpCardMsg, msgEvent.MsgId, "", "")
			//go utils.SendMessage(10, msgEvent.TargetId, InviteCard, "", "", "")

		}
		if msgEvent.Content == "/切歌" {
			go utils.SendMessage(1, msgEvent.TargetId, "功能已经加入待开发队列", msgEvent.MsgId, "", "")
		}
		if strings.HasPrefix(msgEvent.Content, "/playlist") {
			re := regexp.MustCompile(`/\w+ (\d+)`)

			match := re.FindStringSubmatch(msgEvent.Content)
			listId := ""
			if len(match) > 1 {
				listId = match[1]
				go song.GetListAllSongs(listId, msgEvent.GuildID, msgEvent.TargetId, msgEvent.AuthorId, msgEvent.Author.Username)
			} else {
				fmt.Println("播放列表id获取有误！")
			}
		}

		if msgEvent.Content == "/登录" {
			//获取登陆api
			//判断数据是否为空
			go utils.SendMessage(1, msgEvent.TargetId, "二维码登陆，功能已经加入开发队列", msgEvent.MsgId, "", "")
			//存储当前服务器的登陆状态
		}

		//当前bot的状态 播放音乐？当前播放的进度条？下一首预告？
		//当前bot的状态 播放音乐？当前播放的进度条？下一首预告？结束

		//当前bot的播放list 播放列表，超过50条不可添加 按序号排列 可以输入/删除 [2]
		if strings.HasPrefix(msgEvent.Content, "/link") {
			linkData := regexp.MustCompile(`http://\w+.\w+/\w+`)
			getlink := linkData.FindString(msgEvent.Content)
			utils.SendMessage(1, msgEvent.TargetId, getlink, "", "", "")
		}
		//处理网易云音乐
		if strings.HasPrefix(msgEvent.KMarkdown.RawContent, "/网易") {
			re := regexp.MustCompile(`/网易 (.*)`)
			match := re.FindStringSubmatch(msgEvent.KMarkdown.RawContent)
			receiveSongName := ""
			if len(match) > 0 {
				receiveSongName = match[1]
			} else {
				utils.SendMessage(1, msgEvent.TargetId, "客官，关键词有误", "", "", "")
				return err
			}

			//判断用户发送的是手机版链接还是pc端
			//linkData := regexp.MustCompile(`http://\w+.\w+/\w+`)
			//getlink := linkData.FindString(receiveSongName)

			//判断用户是否在语音内
			//获取当前点歌的歌曲id
			songId, songName, songSinger, songPic, err := song.Search(receiveSongName)
			if err != nil {
				return err
			}
			//获取当前点歌的歌曲id结束
			//获取歌曲详情
			songInfo, err := song.MusicInfo(songId)
			fmt.Println("403335371，获取到的歌曲详情：", "歌曲ID:", songId, songInfo, "歌曲名:", songName, "歌手:", songSinger, "专辑图片:", songPic)
			//获取歌曲详情结束
			songid := fmt.Sprintf("%d", songId)

			//将获取歌曲播放url
			//将获取歌曲播放url结束
			//获取当前用户所在的服务

			//创建播放列表数据库
			err = conf.DB.AutoMigrate(&model.Playlist{}, &model.Song{})
			if err != nil {
				return err
			}
			//创建数据库播放列表结束
			//将歌曲信息传入channel
			var playlist model.Playlist
			err = conf.DB.First(&playlist, msgEvent.GuildID).Error
			if err != nil {
				conf.DB.Create(&model.Playlist{ID: msgEvent.GuildID, Songs: []model.Song{{
					SongId:     songid,
					SongName:   songName,
					SongSinger: songSinger,
					CoverUrl:   songPic,
					UserId:     msgEvent.AuthorId,
					UserName:   msgEvent.Author.Username,
				}}})
			} else {
				conf.DB.Create(&model.Song{SongId: songid,
					SongName:   songName,
					CoverUrl:   songPic,
					UserName:   msgEvent.Author.Username,
					SongSinger: songSinger,
					UserId:     msgEvent.AuthorId,
					PlaylistID: msgEvent.GuildID,
				})
			}
			cid, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
			if err != nil {
				return err
			} else if cid == "" {
				utils.SendMessage(1, msgEvent.TargetId, "当前您未处在任何语音频道中！！！", msgEvent.MsgId, "", "")
			} else {
				kook.Play(msgEvent.GuildID, cid, msgEvent.AuthorId)

				MusicCard := model.CardMessageCard{
					Theme: model.CardThemeWarning,
					Size:  model.CardSizeLg,
					Modules: []interface{}{
						&model.CardMessageHeader{Text: model.CardMessageElementText{
							Content: "已将《" + songName + "》添加至列表",
							Emoji:   true,
						}},
						model.CardMessageDivider{},
						&model.CardMessageSection{
							Mode: model.CardMessageSectionModeLeft,
							Text: &model.CardMessageElementKMarkdown{Content: "> " + songName + "\n" + songSinger},
							Accessory: &model.CardMessageElementImage{
								Src:    songPic,
								Alt:    "歌曲专辑图片",
								Size:   "lg",
								Circle: true,
							},
						},
						model.CardMessageDivider{},
						&model.CardMessageContext{
							model.CardMessageElementImage{
								Src: "https://img.kookapp.cn/assets/2023-07/aYf8cNg1hC05k05k.png",
							},
							model.CardMessageElementKMarkdown{Content: "[网易云](https://music.163.com/#/song?id=" + songid + ")"},
						},
					},
				}
				SongCard, err := model.CardMessage{&MusicCard}.BuildMessage()
				if err != nil {
					log.Error(err)
				}
				utils.SendMessage(10, msgEvent.TargetId, SongCard, msgEvent.MsgId, "", "")
			}

			//添加音乐并自动创建播放列表
			//将歌曲添加至频道列表结束

		}
		if msgEvent.Content == "/重连" {
			//获取登陆api
			//判断数据是否为空
			cid, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
			if err != nil {
				return err
			} else if cid == "" {
				utils.SendMessage(1, msgEvent.TargetId, "当前您未处在任何语音频道中！！！", msgEvent.MsgId, "", "")
			} else {
				//player := kook.VoiceInstance{
				//	Token:     conf.Token,
				//	ChannelId: cid,
				//}
				go utils.SendMessage(1, msgEvent.TargetId, "正在执行", msgEvent.MsgId, "", "")
				//存储当前服务器的登陆状态
				err := os.Remove("stream" + cid)
				if err != nil {
					return err
				}
				kook.Status.Delete(cid)
				kook.Play(msgEvent.GuildID, cid, msgEvent.AuthorId)
			}
		}
		//处理网易云音乐结束
		//列出播放列表
		if msgEvent.Content == "/列表" {

			//	查询当前频道的播放列表
			err := kook.PlayForList(msgEvent.GuildID, msgEvent.TargetId)
			if err != nil {
				return err
			}
			//	查询当前频道的播放列表结束
		}

		if msgEvent.Content == "/状态" {

			//err := kook.PlayForList(msgEvent.GuildID, msgEvent.TargetId)
			//if err != nil {
			//	return err
			//}
			//	当前服务器状态
			goinfo := runtime.NumGoroutine()
			goroutineIfo := fmt.Sprintf("%d", goinfo)
			MemPercent, _ := mem.VirtualMemory()
			MemInfo := fmt.Sprintf("%.2f%%", MemPercent.UsedPercent)
			percent, _ := cpu.Percent(time.Second, false)
			CpuInfo := fmt.Sprintf("%.2f%%", percent[0])
			cardData := model.CardMessageCard{
				Theme: model.CardThemeSecondary,
				Color: "",
				Size:  "lg",
				Modules: []interface{}{
					&model.CardMessageHeader{
						Text: model.CardMessageElementText{
							Content: "🌟Status🌟",
							Emoji:   true,
						},
					},
					&model.CardMessageDivider{},
					&model.CardMessageSection{
						Text: model.CardMessageParagraph{
							Cols: 3,
							Fields: []interface{}{
								model.CardMessageElementKMarkdown{Content: "**CPU占用**\n" + CpuInfo},
								model.CardMessageElementKMarkdown{Content: "**内存占用**\n" + MemInfo},
								model.CardMessageElementKMarkdown{Content: "**线 程 数**\n " + goroutineIfo},
							},
						},
					},
					&model.CardMessageDivider{},
				},
			}
			StatusMsg, _ := model.CardMessage{&cardData}.BuildMessage()
			if err != nil {
				return err
			}
			utils.SendMessage(10, msgEvent.TargetId, StatusMsg, "", "", "")
		}
		//列出播放列表结束
		return nil
	}()
	if err != nil {
		log.WithError(err).Error("频道文本事件处理出错！！！")
	}

	return nil
}

type PlayMusicHandler struct {
}

func (play *PlayMusicHandler) Handle(e event.Event) error {
	log.WithField("event", e).Info("音乐播放事件···")
	go func() {
		//	判断服务器id
		var playlist model.Playlist
		conf.DB.Preload("Songs").Find(&playlist)
		if len(playlist.Songs) > 0 {
			var wg sync.WaitGroup
			wg.Add(len(playlist.Songs))
			for _, item := range playlist.Songs {
				channelId, _ := kook.GetChannelId(item.PlaylistID, item.UserId)
				gatewayUrl := kookvoice.GetGatewayUrl(conf.Token, channelId)
				songUrl, times := song.GetMusicUrl(item.SongId)
				connect, rtpUrl := kookvoice.InitWebsocketClient(gatewayUrl)
				defer connect.Close()
				go kookvoice.KeepWebsocketClientAlive(connect)
				go kookvoice.KeepRecieveMessage(connect)
				conf.DB.Debug().Where("id=?", item.ID).Delete(&playlist.Songs)
				kookvoice.StreamAudio(rtpUrl, songUrl)
				wg.Done()
				fmt.Println(times)
			}
			wg.Wait()

		}
	}()
	return nil
}
