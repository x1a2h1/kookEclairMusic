package app

import (
	"botserver/app/kook"
	"botserver/app/model"
	"botserver/app/netease"
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

		if msgEvent.Content == "/帮助" {
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
						Text: model.CardMessageElementKMarkdown{Content: "**(font)网易云(font)[pink]**\n> `/网易 value` or `/wy value` \n`value`为歌单链接、歌曲名、歌曲链接\n(spl)歌单链接默认导入0-50首(spl)\n\n---\n**(font)其他指令(font)[purple]**\n> `/列表`:获取当前服务器播放列表。\n`/切歌`:播放下一首。\n`/退出`:退出并清空播放列表\n\n---"},
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
			go utils.SendMessage(10, msgEvent.TargetId, helpCardMsg, msgEvent.MsgId, "", "")
			//go utils.SendMessage(10, msgEvent.TargetId, InviteCard, "", "", "")

		}
		if msgEvent.Content == "/切歌" {
			//终止goroutine
			_, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
			if err != nil {
				go utils.SendMessage(1, msgEvent.TargetId, "!ok!", msgEvent.MsgId, "", "")
				return err
			} else {
				_, ok := kook.Status.Load(msgEvent.GuildID)
				if ok {
					kook.Mu.Lock()
					cmds, ok := kook.Cmds[msgEvent.GuildID]
					kook.Mu.Unlock()
					if ok && cmds != nil && cmds.Process != nil {
						err := cmds.Process.Kill()
						if err != nil {
							log.Printf("Error while killing ffmpeg for %s: %v", msgEvent.GuildID, err)
						}
						go utils.SendMessage(1, msgEvent.TargetId, "ok!", msgEvent.MsgId, "", "")
					} else {
						log.Error("当前服务器没有ffmpeg进程")
					}
					//	如果服务没有正在播放，判断当前歌单中是否存在歌曲
					//var songList model.Song
					//	如果
				} else {
					log.Error("当前频道没有在播放")
				}
			}
			//kook.Player(msgEvent.GuildID, cid, msgEvent.AuthorId)
			//根据服务器id 当前切歌用户id判断，当前服务器歌单中是否有该用户的歌曲，判断当前服务器播放拼
			//判断当前服务器是否在播放 没有则返回，如果正在播放 终止当前的for

			//go utils.SendMessage(1, msgEvent.TargetId, "功能已经加入待开发队列", msgEvent.MsgId, "", "")
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
		//处理网易云音乐
		if strings.HasPrefix(msgEvent.Content, "/网易") || strings.HasPrefix(msgEvent.Content, "/wy") {
			re := regexp.MustCompile(`/(网易|wy) (.*)`)
			match := re.FindStringSubmatch(msgEvent.Content)
			receiveSongName := ""
			songId := ""
			if len(match) > 0 {
				link := regexp.MustCompile(`https?://`)
				if link.MatchString(match[2]) {
					//将链接进行处理获取id
					linkPattern := regexp.MustCompile(`(https?://[^\s\]]+)`)
					linkMatch := linkPattern.FindStringSubmatch(match[2])
					fullLink := linkMatch[1]
					fmt.Println("link的值为：", fullLink)
					params, id, err := netease.HandleLink(fullLink)
					if err != nil {
						log.Error("处理链接有误！", err)
					}
					fmt.Println("当前获取到的值为："+params, "id为：", id)
					switch params {
					case "playlist":
						cid, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
						if err != nil || cid == "" {
							utils.SendMessage(1, msgEvent.TargetId, "获取播放频道失败或您未处在任何语音频道！！", "", "", "")
							break
						}
						err = conf.DB.AutoMigrate(&model.Playlist{}, &model.Song{})
						if err != nil {
							return err
							break
						}
						go song.GetListAllSongs(id, msgEvent.GuildID, msgEvent.TargetId, msgEvent.AuthorId, cid, msgEvent.Author.Username)
						time.Sleep(time.Second)
						_, ok := kook.Status.Load(msgEvent.GuildID)
						if !ok {
							err := kook.Play(msgEvent.GuildID, cid, msgEvent.AuthorId)
							if err != nil {
								return err
							}
						}
						return nil
					case "song":
						songId = id
					default:
						utils.SendMessage(1, msgEvent.TargetId, "链接有误！", "", "", "")
					}
				} else {
					receiveSongName = match[2]
					id, err := song.Search(receiveSongName)
					if err != nil {
						return err
					}
					songid := fmt.Sprintf("%d", id)
					songId = songid
				}
			} else {
				utils.SendMessage(1, msgEvent.TargetId, "客官，关键词有误", "", "", "")
				return err
			}
			//判断用户是否在语音内

			//获取当前点歌的歌曲id
			//songId, songName, songSinger, songPic, err := song.Search(receiveSongName)

			//获取当前点歌的歌曲id结束
			//获取歌曲详情
			songName, songSinger, songPic, dt, err := song.MusicInfo(songId)
			fmt.Println("403335371，获取到的歌曲详情：", "歌曲ID:", songId, "歌曲名:", songName, "歌手:", songSinger, "专辑图片:", songPic, "歌曲时长", dt)
			//获取歌曲详情结束
			//songid := fmt.Sprintf("%d", songId)

			//将获取歌曲播放url
			//将获取歌曲播放url结束
			//获取当前用户所在的服务

			//创建播放列表数据库
			err = conf.DB.AutoMigrate(&model.Playlist{}, &model.Song{})
			if err != nil {
				return err
			}
			//创建数据库播放列表结束

			cid, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
			if err != nil {
				return err
			} else if cid == "" {
				utils.SendMessage(1, msgEvent.TargetId, "当前您未处在任何语音频道中！！！", msgEvent.MsgId, "", "")
			} else {
				//将歌曲写入数据库
				var playlist model.Playlist
				err = conf.DB.First(&playlist, msgEvent.GuildID).Error
				if err != nil {
					conf.DB.Create(&model.Playlist{ID: msgEvent.GuildID, Songs: []model.Song{{
						SongId:     songId,
						SongName:   songName,
						SongSinger: songSinger,
						CoverUrl:   songPic,
						UserId:     msgEvent.AuthorId,
						UserName:   msgEvent.Author.Username,
					}}})
				} else {
					conf.DB.Create(&model.Song{SongId: songId,
						SongName:   songName,
						CoverUrl:   songPic,
						UserName:   msgEvent.Author.Username,
						SongSinger: songSinger,
						UserId:     msgEvent.AuthorId,
						PlaylistID: msgEvent.GuildID,
					})
				}
				//进行播放
				kook.Play(msgEvent.GuildID, cid, msgEvent.AuthorId)
				//发送卡片
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
							model.CardMessageElementKMarkdown{Content: "[网易云](https://music.163.com/#/song?id=" + songId + ")"},
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
		if msgEvent.Content == "/退出" || msgEvent.Content == "/quit" {
			//判断当前用户是否在语音频道当中，在进行判断当前频道是否在播放
			cid, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
			if err != nil {
				return err
			} else if cid == "" {
				utils.SendMessage(1, msgEvent.TargetId, "不可操作！！！", msgEvent.MsgId, "", "")
				return err
			}
			_, ok := kook.Status.Load(msgEvent.GuildID)
			if ok {
				//	删数据并终止进程
				conf.DB.Where("playlist_id = ?", msgEvent.GuildID).Delete(&model.Song{})
				kook.Mu.Lock()
				cmds, ok := kook.Cmds[msgEvent.GuildID]
				kook.Mu.Unlock()
				if ok && cmds != nil && cmds.Process != nil {
					err := cmds.Process.Kill()
					if err != nil {
						log.Printf("Error while killing ffmpeg for %s: %v", msgEvent.GuildID, err)
					}
					go utils.SendMessage(1, msgEvent.TargetId, "ok!", msgEvent.MsgId, "", "")
				}
			}
		}
		if msgEvent.Content == "/重连" {
			//判断当前服务器status是否为true
			//如果ok，关闭当前服务器goroutine
			// ！ok 判断当前服务器列表是否为空，不为空进行重连
			cid, err := kook.GetChannelId(msgEvent.GuildID, msgEvent.AuthorId)
			if err != nil {
				return err
			} else if cid == "" {
				utils.SendMessage(1, msgEvent.TargetId, "当前您未处在任何语音频道中！！！", msgEvent.MsgId, "", "")
			} else {
				_, ok := kook.Status.Load(msgEvent.GuildID)
				if !ok {
					err := kook.Play(msgEvent.GuildID, cid, msgEvent.AuthorId)
					if err != nil {
						return err
					}
				} else {

				}
				go utils.SendMessage(1, msgEvent.TargetId, "ok!", msgEvent.MsgId, "", "")
				//存储当前服务器的登陆状态
				//err := os.Remove("stream" + cid)
				//if err != nil {
				//	return err
				//}
				//kook.Status.Delete(msgEvent.GuildID)
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
			C := 0
			kook.Status.Range(func(key, value any) bool {
				C++
				return true
			})
			CPlaying := ""
			if C == 0 {
				CPlaying = "无"
			} else {
				CPlaying = fmt.Sprintf("(font)%d(font)[warning] 个频道", C)
			}
			totalPlay := fmt.Sprintf("(font)%d(font)[warning]首", kook.TotalPlay)
			var musicTotal int64
			var songs model.Song
			conf.DB.Debug().Model(&songs).Count(&musicTotal)
			total := fmt.Sprintf("%d", musicTotal)
			goinfo := runtime.NumGoroutine()
			goroutineIfo := fmt.Sprintf("%d", goinfo)
			MemPercent, _ := mem.VirtualMemory()
			MemInfo := ""
			if MemPercent.UsedPercent < 40 {
				MemInfo = fmt.Sprintf("(font)%.2f%%(font)[success]", MemPercent.UsedPercent)
			} else if MemPercent.UsedPercent < 60 {
				MemInfo = fmt.Sprintf("(font)%.2f%%(font)[primary]", MemPercent.UsedPercent)
			} else if MemPercent.UsedPercent < 80 {
				MemInfo = fmt.Sprintf("(font)%.2f%%(font)[warning]", MemPercent.UsedPercent)
			} else {
				MemInfo = fmt.Sprintf("(font)%.2f%%(font)[danger]", MemPercent.UsedPercent)
			}
			percent, _ := cpu.Percent(time.Second, false)
			CpuInfo := ""
			if percent[0] < 40 {
				CpuInfo = fmt.Sprintf("(font)%.2f%%(font)[success]", percent[0])
			} else if percent[0] < 60 {
				CpuInfo = fmt.Sprintf("(font)%.2f%%(font)[primary]", percent[0])
			} else if percent[0] < 80 {
				CpuInfo = fmt.Sprintf("(font)%.2f%%(font)[warning]", percent[0])
			} else {
				CpuInfo = fmt.Sprintf("(font)%.2f%%(font)[danger]", percent[0])
			}

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
								model.CardMessageElementKMarkdown{Content: "**其  他**\n (font)" + goroutineIfo + "(font)[purple]"},
								model.CardMessageElementKMarkdown{Content: "**待 播 放**\n (font)" + total + "(font)[warning]首"},
								model.CardMessageElementKMarkdown{Content: "**正在服务**\n " + CPlaying},
								model.CardMessageElementKMarkdown{Content: "**总 播 放**\n " + totalPlay},
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
