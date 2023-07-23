package song

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
)

func init() {

}

type MusicUrl struct {
	Data []struct {
		Url string `json:"url"`
	} `json:"data"`
}

type Response struct {
	Result struct {
		SearchQcReminder interface{} `json:"searchQcReminder"`
		Songs            []struct {
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
			Rt   interface{}   `json:"rt"`
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
			Sq                   interface{}   `json:"sq"`
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
			Single               int           `json:"single"`
			NoCopyrightRcmd      interface{}   `json:"noCopyrightRcmd"`
			Rtype                int           `json:"rtype"`
			Rurl                 interface{}   `json:"rurl"`
			Mst                  int           `json:"mst"`
			Cp                   int           `json:"cp"`
			Mv                   int           `json:"mv"`
			PublishTime          int64         `json:"publishTime"`
			Privilege            struct {
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
					ResConsumable      bool        `json:"resConsumable"`
					UserConsumable     bool        `json:"userConsumable"`
					ListenType         interface{} `json:"listenType"`
					CannotListenReason interface{} `json:"cannotListenReason"`
				} `json:"freeTrialPrivilege"`
				RightSource    int `json:"rightSource"`
				ChargeInfoList []struct {
					Rate          int         `json:"rate"`
					ChargeUrl     interface{} `json:"chargeUrl"`
					ChargeMessage interface{} `json:"chargeMessage"`
					ChargeType    int         `json:"chargeType"`
				} `json:"chargeInfoList"`
			} `json:"privilege"`
		} `json:"songs"`
		SongCount int `json:"songCount"`
	} `json:"result"`
	Code int `json:"code"`
}

func Search(keywords string) (int, string, string, error) {
	//	通过关键词 搜索歌曲 将返回body转化成json可读形式，将lenght：0 的id传给获取播放地址 在返回进行播放伴奏
	//	如果/网易 明明就 周杰伦  搜索多匹配 歌曲名 歌手
	keywords = url.QueryEscape(keywords)

	resp, err := http.Get("http://192.168.110.69:3000/cloudsearch?keywords=" + keywords + "&limit=1")
	if err != nil {
		return 0, "", "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", "", err
	}
	var result Response
	err = json.Unmarshal(body, &result)
	if err != nil {
		return 0, "", "", err
	}
	id := result.Result.Songs[0].Id
	name := result.Result.Songs[0].Name
	pic := result.Result.Songs[0].Al.PicUrl
	return id, name, pic, nil
}

func GetMusicUrl(id string) string {
	resp, err := http.Get("http://192.168.110.69:3000/song/url/v1?id=" + id + "&level=exhigh")
	if err != nil {
		log.Error("403335371获取音乐url出现错误！", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	var res MusicUrl
	err = json.Unmarshal(body, &res)
	if err != nil {
		return ""
	}
	songurl := res.Data[0].Url

	return songurl
}
