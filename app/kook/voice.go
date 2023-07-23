package kook

import (
	model "botserver/app/model"
	"botserver/conf"
	"botserver/pkg/untils"
	"fmt"
)

func PlayForList(gid string, targerId string) error {
	//	根据播放列表播放
	var Current = make(map[string]string)
	Current["songName"] = "陈宇晖牛逼"
	Current["songUrl"] = "https://c-ssl.dtstatic.com/uploads/blog/202207/09/20220709150824_97667.thumb.400_0.jpg"
	Current["singer"] = "陈宇晖"
	Current["userName"] = "夏至"
	var playlist model.Playlist
	conf.DB.Preload("Songs").Take(&playlist, gid)

	//var forList []interface{}
	//for _, item := range playlist.Songs {
	//	data :=
	//	forList = append(forList, data)
	//}
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
