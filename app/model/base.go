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

// 写入歌曲信息并创建播放列表
func CreateList(gid string, mid string, name string, user string, cover string) {
}
