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

type SongInfo struct {
	Songs []struct {
		Name string `json:"name"`
		Id   int    `json:"id"`
		Pst  int    `json:"pst"`
		T    int    `json:"t"`
		Ar   []struct {
			Id    int           `json:"id"`
			Name  string        `json:"name"`
			Tns   []interface{} `json:"tns"`
			Alias []interface{} `json:"alias"`
		} `json:"ar"`
		Alia []interface{} `json:"alia"`
		Pop  int           `json:"pop"`
		St   int           `json:"st"`
		Rt   string        `json:"rt"`
		Fee  int           `json:"fee"`
		V    int           `json:"v"`
		Crbt interface{}   `json:"crbt"`
		Cf   string        `json:"cf"`
		Al   struct {
			Id     int           `json:"id"`
			Name   string        `json:"name"`
			PicUrl string        `json:"picUrl"`
			Tns    []interface{} `json:"tns"`
			PicStr string        `json:"pic_str"`
			Pic    int64         `json:"pic"`
		} `json:"al"`
		Dt int `json:"dt"`
		H  struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"h"`
		M struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"m"`
		L struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"l"`
		Sq struct {
			Br   int `json:"br"`
			Fid  int `json:"fid"`
			Size int `json:"size"`
			Vd   int `json:"vd"`
			Sr   int `json:"sr"`
		} `json:"sq"`
		Hr                   interface{}   `json:"hr"`
		A                    interface{}   `json:"a"`
		Cd                   string        `json:"cd"`
		No                   int           `json:"no"`
		RtUrl                interface{}   `json:"rtUrl"`
		Ftype                int           `json:"ftype"`
		RtUrls               []interface{} `json:"rtUrls"`
		DjId                 int           `json:"djId"`
		Copyright            int           `json:"copyright"`
		SId                  int           `json:"s_id"`
		Mark                 int           `json:"mark"`
		OriginCoverType      int           `json:"originCoverType"`
		OriginSongSimpleData interface{}   `json:"originSongSimpleData"`
		TagPicList           interface{}   `json:"tagPicList"`
		ResourceState        bool          `json:"resourceState"`
		Version              int           `json:"version"`
		SongJumpInfo         interface{}   `json:"songJumpInfo"`
		EntertainmentTags    interface{}   `json:"entertainmentTags"`
		AwardTags            interface{}   `json:"awardTags"`
		Single               int           `json:"single"`
		NoCopyrightRcmd      interface{}   `json:"noCopyrightRcmd"`
		Rtype                int           `json:"rtype"`
		Rurl                 interface{}   `json:"rurl"`
		Mst                  int           `json:"mst"`
		Cp                   int           `json:"cp"`
		Mv                   int           `json:"mv"`
		PublishTime          int64         `json:"publishTime"`
	} `json:"songs"`
	Privileges []struct {
		Id                 int         `json:"id"`
		Fee                int         `json:"fee"`
		Payed              int         `json:"payed"`
		St                 int         `json:"st"`
		Pl                 int         `json:"pl"`
		Dl                 int         `json:"dl"`
		Sp                 int         `json:"sp"`
		Cp                 int         `json:"cp"`
		Subp               int         `json:"subp"`
		Cs                 bool        `json:"cs"`
		Maxbr              int         `json:"maxbr"`
		Fl                 int         `json:"fl"`
		Toast              bool        `json:"toast"`
		Flag               int         `json:"flag"`
		PreSell            bool        `json:"preSell"`
		PlayMaxbr          int         `json:"playMaxbr"`
		DownloadMaxbr      int         `json:"downloadMaxbr"`
		MaxBrLevel         string      `json:"maxBrLevel"`
		PlayMaxBrLevel     string      `json:"playMaxBrLevel"`
		DownloadMaxBrLevel string      `json:"downloadMaxBrLevel"`
		PlLevel            string      `json:"plLevel"`
		DlLevel            string      `json:"dlLevel"`
		FlLevel            string      `json:"flLevel"`
		Rscl               interface{} `json:"rscl"`
		FreeTrialPrivilege struct {
			ResConsumable  bool        `json:"resConsumable"`
			UserConsumable bool        `json:"userConsumable"`
			ListenType     interface{} `json:"listenType"`
		} `json:"freeTrialPrivilege"`
		ChargeInfoList []struct {
			Rate          int         `json:"rate"`
			ChargeUrl     interface{} `json:"chargeUrl"`
			ChargeMessage interface{} `json:"chargeMessage"`
			ChargeType    int         `json:"chargeType"`
		} `json:"chargeInfoList"`
	} `json:"privileges"`
	Code int `json:"code"`
}

func MusicInfo(id int) (map[string]interface{}, error) {
	//	获取歌曲详情
	url := fmt.Sprintf(conf.NetEasy+"/song/detail?ids=%d", id)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res SongInfo
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}
	fmt.Println("403335371，获取到的歌曲详情为", res)
	return nil, err
}

// 获取歌单歌曲
func GetPlaylist(id string) {
	//	对所有歌曲写库

}

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

func GetListAllSongs(id string, gid string, targetId string, uid string, uname string) error {
	//获取歌单中所有歌曲
	resp, err := http.Get(conf.NetEasy + "/playlist/detail?id=" + id)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err

	}
	var res ListInfo
	err = json.Unmarshal(body, &res)

	SongTotal := fmt.Sprintf("%d", len(res.Playlist.Tracks))

	//写入数据库
	for _, songItem := range res.Playlist.Tracks {
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
	//cid, err := kook.GetChannelId(gid, uid)
	//kook.Play(gid, "3970172859382929", uid)

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

	for _, item := range res.Playlist.Tracks {
		listCard.AddModule(
			&model.CardMessageSection{
				Mode: "left",
				Text: model.CardMessageElementKMarkdown{
					Content: "> " + item.Name + "\n> " + item.Ar[0].Name,
				},
				//Accessory: &model.CardMessageElementImage{
				//	Src:  item.Al.PicUrl,
				//	Size: "lg",
				//},
			},
		)
	}
	sendMsg, err := model.CardMessage{&listCard}.BuildMessage()
	if err != nil {
		return err
	}
	err = utils.SendMessage(10, targetId, sendMsg, "", "", "")
	if err != nil {
		return err
	}
	fmt.Println(gid, "当前歌单的数据为", res.Playlist)
	return err
}
