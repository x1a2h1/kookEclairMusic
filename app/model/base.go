package model

type Playlist struct {
	ID    string `gorm:"size:16"`
	Songs []Song
}

type Song struct {
	SongId     string //歌曲id
	SongName   string //歌曲名
	CoverUrl   string //歌曲头像
	UserName   string //点歌用户名
	PlaylistID string //当前歌曲属于哪个服务器的ID
}

// 写入歌曲信息并创建播放列表
func CreateList(gid string, mid string, name string, user string, cover string) {

}
