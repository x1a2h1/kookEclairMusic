package model

type Playlist struct {
	ID    string `gorm:"size:16"`
	Songs []Song
}

type ListInfo struct {
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

type Song struct {
	ID         int    //自增id
	SongId     string //歌曲id
	SongName   string //歌曲名
	SongSinger string //歌手名
	CoverUrl   string //专辑图片
	UserName   string //点歌用户名
	UserId     string //点歌用户id
	PlaylistID string //当前歌曲属于哪个服务器的ID
}

// channel中的歌曲信息
type MusicList struct {
	Guild    string //服务器id
	ChanId   string //用户当前所在的语音频道
	SongId   string //当前歌曲id
	SongName string //歌曲名
	MusicUrl string //歌曲播放地址
	UserName string //点歌用户
	CoverUrl string //歌曲图片&专辑图片
	Duration int
}
type GatewayResp struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Data    gatewayRespData `json:"data"`
}
type gatewayRespData struct {
	GatewayUrl string `json:"gateway_url"`
}
type FirstShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    firstShakeReqData `json:"data"`
}
type firstShakeReqData struct {
}

type SecondShakeReq struct {
	Request bool               `json:"request"`
	Id      int                `json:"id"`
	Method  string             `json:"method"`
	Data    SecondShakeReqData `json:"data"`
}

type SecondShakeReqData struct {
	DisplayName string `json:"displayName"`
}
type BaseShakeReq struct {
	Request bool   `json:"request"`
	Id      int    `json:"id"`
	Method  string `json:"method"`
}

type ThirdShakeReq struct {
	Request bool              `json:"request"`
	Id      int               `json:"id"`
	Method  string            `json:"method"`
	Data    ThirdShakeReqData `json:"data"`
}

type ThirdShakeReqData struct {
	Comedia bool   `json:"comedia"`
	RtcpMux bool   `json:"rtcpMux"`
	Type    string `json:"type"`
}

type ThirdShakeResp struct {
	Response bool `json:"response"`
	Id       int  `json:"id"`
	Ok       bool `json:"ok"`
	Data     thirdShakeRespData
}

type thirdShakeRespData struct {
	Id       string `json:"id"`
	Ip       string `json:"ip"`
	Port     int    `json:"port"`
	RtcpPort int    `json:"rtcpPort"`
}

type FourthShakeReq struct {
	Request bool               `json:"request"`
	Id      int                `json:"id"`
	Method  string             `json:"method"`
	Data    FourthShakeReqData `json:"data"`
}

type FourthShakeReqData struct {
	AppData       AppData       `json:"appData"`
	Kind          string        `json:"kind"`
	PeerId        string        `json:"peerId"`
	RtpParameters RtpParameters `json:"rtpParameters"`
	TransportId   string        `json:"transportId"`
}

type AppData struct {
}

type RtpParameters struct {
	Codecs    []Codec    `json:"codecs"`
	Encodings []Encoding `json:"encodings"`
}

type Codec struct {
	Channels    int        `json:"channels"`
	ClockRate   int        `json:"clockRate"`
	MimeType    string     `json:"mimeType"`
	Parameters  Parameters `json:"parameters"`
	PayloadType int        `json:"payloadType"`
}

type Parameters struct {
	SpropStereo int `json:"sprop-stereo"`
}

type Encoding struct {
	Ssrc int `json:"ssrc"`
}

// 写入歌曲信息并创建播放列表
func CreateList(gid string, mid string, name string, user string, cover string) {
}
