package app

import (
	"botserver/app/model"
	"botserver/app/song"
	"botserver/conf"
	"botserver/pkg/untils"
	"errors"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/gookit/event"
	"github.com/idodo/golang-bot/kaihela/api/base"
	event2 "github.com/idodo/golang-bot/kaihela/api/base/event"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
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
		//当前bot的状态 播放音乐？当前播放的进度条？下一首预告？
		//当前bot的状态 播放音乐？当前播放的进度条？下一首预告？结束

		//当前bot的播放list 播放列表，超过50条不可添加 按序号排列 可以输入/删除 [2]

		//处理网易云音乐
		if strings.HasPrefix(msgEvent.KMarkdown.RawContent, "/网易") {
			re := regexp.MustCompile(`/网易\s+(\S+)`)
			match := re.FindStringSubmatch(msgEvent.KMarkdown.RawContent)
			var receiveSongName string = match[1]
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
		}
		//处理网易云音乐结束
		return nil
	}()
	if err != nil {
		log.WithError(err).Error("频道文本事件处理出错！！！")

	}
	return nil
}
