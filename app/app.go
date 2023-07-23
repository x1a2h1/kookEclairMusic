package app

import (
	"botserver/app/kook"
	"botserver/app/model"
	_ "botserver/app/redis"
	"botserver/app/song"
	"botserver/conf"
	"botserver/pkg/untils"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	"github.com/kaiheila/golang-bot/api/base"
	event2 "github.com/kaiheila/golang-bot/api/base/event"
	"github.com/kaiheila/golang-bot/api/helper"
	"github.com/shuyangzhang/kookvoice"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var dataMap map[string]interface{}

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
		//go Guild.HandleServerEvents(msgEvent.GuildID)
		dataMap = make(map[string]interface{})
		//开启多线程结束

		helpCard := model.CardMessageCard{
			Theme: model.CardThemeDanger,
			Size:  model.CardSizeLg,
			Modules: []interface{}{
				&model.CardMessageHeader{Text: model.CardMessageElementText{
					Content: "🌈帮助菜单 & help menu",
					Emoji:   true,
				}},
				&model.CardMessageDivider{},
				&model.CardMessageSection{
					Text: model.CardMessageParagraph{
						Cols: 3,
						Fields: []interface{}{
							model.CardMessageElementKMarkdown{Content: "**指令**\n(font)/网易 { 歌曲名 } (font)[error]\n(font)/QQ(font)[success]"},
							model.CardMessageElementKMarkdown{Content: "**功能**\n(font)播放网易云音乐(font)[success]\n待完善"},
							model.CardMessageElementKMarkdown{Content: "**示例**\n/网易 乐鼓 (dj版)\n待完善"},
						},
					},
				},
				&model.CardMessageDivider{},
				&model.CardMessageContext{
					&model.CardMessageElementText{Content: "当前频道id：" + msgEvent.TargetId + "\n"},
					&model.CardMessageElementText{Content: "当前频道名：" + msgEvent.ChannelName + "\n"},
					&model.CardMessageElementText{Content: "当前频道服务器ID：" + msgEvent.GuildID + "\n"},
					&model.CardMessageElementText{Content: "当前频道服务器ID：" + msgEvent.Nonce + "\n"},
				},
				&model.CardMessageSection{
					Text: model.CardMessageElementKMarkdown{Content: "Version:" + "`" + conf.Version + "`"},
				},
			},
		}

		helpCardMsg, err := model.CardMessage{&helpCard}.BuildMessage()
		if err != nil {
			log.Error("编译信息时出错！", err)
		}
		if msgEvent.Content == "/帮助" {
			untils.SendMessage(10, msgEvent.TargetId, helpCardMsg, msgEvent.MsgId, "", "")
		}
		if msgEvent.Content == "/登录" {
			//获取登陆api
			//判断数据是否为空
			untils.SendMessage(1, msgEvent.TargetId, "二维码登陆，功能待完善", msgEvent.MsgId, "", "")
			//存储当前服务器的登陆状态
		}
		//当前bot的状态 播放音乐？当前播放的进度条？下一首预告？
		//当前bot的状态 播放音乐？当前播放的进度条？下一首预告？结束

		//当前bot的播放list 播放列表，超过50条不可添加 按序号排列 可以输入/删除 [2]

		//处理网易云音乐
		if strings.HasPrefix(msgEvent.KMarkdown.RawContent, "/网易") {
			re := regexp.MustCompile(`/网易\s+(\S+)`)
			match := re.FindStringSubmatch(msgEvent.KMarkdown.RawContent)
			var receiveSongName = match[1]
			//判断用户是否在语音内
			//获取当前点歌的歌曲id
			songId, songName, songPic, err := song.Search(receiveSongName)
			if err != nil {
				return err
			}
			fmt.Println("403335371,这", songId)
			//获取当前点歌的歌曲id结束

			//获取歌曲详情
			songInfo, err := song.MusicInfo(songId)
			fmt.Println("403335371，获取到的歌曲详情：", songInfo, "歌曲名", songName, "歌曲图片", songPic)
			//获取歌曲详情结束
			songid := fmt.Sprintf("%d", songId)

			//将获取歌曲播放url
			musicUrl := song.GetMusicUrl(songid)
			//将获取歌曲播放url结束
			fmt.Println("403335371请求获取歌曲url地址", musicUrl)
			//获取当前用户所在的服务
			Client := helper.NewApiHelper("/v3/channel-user/get-joined-channel", conf.Token, conf.BaseUrl, "", "")
			Client.SetQuery(map[string]string{
				"guild_id": msgEvent.GuildID,
				"user_id":  msgEvent.AuthorId,
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
			//将歌曲添加至频道列表
			//err = conf.DB.AutoMigrate(&model.Playlist{}, &model.Song{})
			//if err != nil {
			//	return err
			//}
			//创建播放列表
			var playlist model.Playlist
			err = conf.DB.First(&playlist, msgEvent.GuildID).Error
			if err != nil {
				conf.DB.Create(&model.Playlist{ID: msgEvent.GuildID, Songs: []model.Song{{
					SongId:   songid,
					SongName: songName,
					CoverUrl: songPic,
					UserName: msgEvent.Author.Username,
				}}})
			} else {
				conf.DB.Create(&model.Song{SongId: songid,
					SongName:   songName,
					CoverUrl:   songPic,
					UserName:   msgEvent.Author.Username,
					PlaylistID: msgEvent.GuildID,
				})
			}
			//添加音乐并自动创建播放列表
			//将歌曲添加至频道列表结束
			MusicCard := model.CardMessageCard{
				Theme: model.CardThemePrimary,
				Size:  model.CardSizeLg,
			}
			cardHeader := &model.CardMessageHeader{Text: model.CardMessageElementText{
				Content: "已将" + songName + "添加至列表",
			}}
			MusicCardSection := &model.CardMessageSection{
				Text: model.CardMessageElementText{
					Content: songName,
				},
				Accessory: model.CardMessageElementImage{
					Src:    songPic,
					Size:   "lg",
					Circle: true,
				},
			}
			MusicCard.AddModule(cardHeader, MusicCardSection)
			msg := model.CardMessage{&MusicCard}
			content, _ := msg.BuildMessage()
			untils.SendMessage(10, msgEvent.TargetId, content, msgEvent.MsgId, "", "")
			go kookvoice.Play(conf.Token, res.Data.Items[0].Id, musicUrl)
		}
		//处理网易云音乐结束
		//列出播放列表
		if msgEvent.Content == "/列表" {
			kook.PlayForList(msgEvent.GuildID, msgEvent.TargetId)

			//	查询当前频道的播放列表
			//	查询当前频道的播放列表结束
			//songList := d
			//data := model.CardMessageCard{
			//	Theme: model.CardThemeSuccess,
			//	Modules: []interface{}{
			//		&model.CardMessageHeader{Text: model.CardMessageElementText{
			//			Content: "播放列表",
			//			Emoji:   false,
			//		}},
			//		&model.CardMessageDivider{},
			//		&model.CardMessageSection{
			//			Mode: "right",
			//			Text: model.CardMessageElementKMarkdown{Content: "> ** **\n> **音乐名**\n> **歌手**\n> ** **"},
			//			Accessory: &model.CardMessageElementImage{
			//				Src:    "https://c-ssl.dtstatic.com/uploads/blog/202207/09/20220709150824_97667.thumb.400_0.jpg",
			//				Size:   "sm",
			//				Circle: true,
			//			},
			//		},
			//		&model.CardMessageDivider{},
			//	},
			//}
			//listMsg, err := model.CardMessage{&data}.BuildMessage()
			//if err != nil {
			//	log.Error("编译信息时出错！", err)
			//	return err
			//}
			//err = untils.SendMessage(10, msgEvent.TargetId, listMsg, "", "", "")
			//if err != nil {
			//	return err
			//}
		}
		//列出播放列表结束
		return nil
	}()
	if err != nil {
		log.WithError(err).Error("频道文本事件处理出错！！！")
	}
	return nil
}
