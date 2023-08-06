package song

import (
	"botserver/app/model"
	"botserver/conf"
	"botserver/pkg/utils"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Info struct {
	Songs []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
		Ar   []struct {
			Name string `json:"name"`
		} `json:"ar"`
		Al struct {
			PicUrl string `json:"picUrl"`
		} `json:"al"`
		Dt int `json:"dt"`
	} `json:"songs"`
	Code int `json:"code"`
}

func MusicInfo(id string) (string, string, string, int, error) {
	//	获取歌曲详情
	//返回歌曲歌曲名、歌手、专辑图片、音频时长
	url := fmt.Sprintf(conf.NetEasy + "/song/detail?ids=" + id)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", "", 0, err
	}
	var res Info
	err = json.Unmarshal(body, &res)
	if err != nil {
		return "", "", "", 0, err
	}
	Singer := ""
	for i, name := range res.Songs[0].Ar {
		if i > 0 {
			Singer += " / "
		}
		Singer += name.Name
		fmt.Println("获取到的歌手名", i, "为：", name.Name)
	}
	fmt.Println("403335371，获取到的歌曲详情为", res)
	return res.Songs[0].Name, Singer, res.Songs[0].Al.PicUrl, res.Songs[0].Dt, err
}

// 获取歌单歌曲

type ListInfo struct {
	//Playlist struct {
	//	Tracks []struct {
	//		Name string `json:"name"`
	//		Id   string `json:"id"`
	//		Ar   []struct {
	//			Name string `json:"name"`
	//		} `json:"ar"`
	//		Al struct {
	//			PicUrl string `json:"picUrl"`
	//		} `json:"al"`
	//	} `json:"tracks"`
	//} `json:"playlist"`

	Playlist struct {
		Id     int64 `json:"id"`
		Tracks []struct {
			Name string `json:"name"`
			Id   int    `json:"id"`
			Ar   []struct {
				Id   int    `json:"id"`
				Name string `json:"name"`
			} `json:"ar"`
			Al struct {
				PicUrl string `json:"picUrl"`
			} `json:"al"`
		} `json:"tracks"`
	} `json:"playlist"`
}

func GetListAllSongs(id string, gid string, targetId string, uid string, chanId string, uname string) error {
	//获取歌单中所有歌曲
	resp, err := http.Get(conf.NetEasy + "/playlist/track/all?id=" + id + "&limit=50")
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	var res model.ListInfo
	err = json.Unmarshal(body, &res)

	SongTotal := fmt.Sprintf("%d", len(res.Songs))

	//写入数据库
	for _, songItem := range res.Songs {
		songId := fmt.Sprintf("%d", songItem.Id)
		conf.DB.Create(&model.Song{
			SongId:     songId,
			SongName:   songItem.Name,
			SongSinger: songItem.Ar[0].Name,
			CoverUrl:   songItem.Al.PicUrl,
			UserName:   uname,
			UserId:     uid,
			PlaylistID: gid,
		})
	}
	//发送卡片
	listCard := model.CardMessageCard{
		Theme: model.CardThemeSuccess,
		Modules: []interface{}{
			&model.CardMessageHeader{Text: model.CardMessageElementText{
				Content: "成功导入" + SongTotal + "首歌曲",
				Emoji:   false,
			}},
			&model.CardMessageDivider{},
		},
	}

	for index, item := range res.Songs {
		if index > 10 {
			break
		}
		listCard.AddModule(
			&model.CardMessageSection{
				Mode: "left",
				Text: model.CardMessageElementKMarkdown{
					Content: "> " + item.Name + "\n> " + item.Ar[0].Name,
				},
				Accessory: &model.CardMessageElementImage{
					Src:  item.Al.PicUrl,
					Size: "sm",
				},
			},
			&model.CardMessageDivider{},
		)
	}
	listCard.AddModule(
		&model.CardMessageSection{Text: model.CardMessageElementText{Content: "当前卡片仅展示10条数据"}},
	)

	sendMsg, err := model.CardMessage{&listCard}.BuildMessage()
	if err != nil {
		return err
	}
	err = utils.SendMessage(10, targetId, sendMsg, "", "", "")
	if err != nil {
		return err
	}
	fmt.Println(gid, "当前歌单的数据为", res.Songs)
	return err
}
