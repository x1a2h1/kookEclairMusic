package model

type Playlist struct {
	ID    string `gorm:"size:36"`
	Songs []Song `gorm:"foreignKey:PlaylistID"`
}

type Song struct {
	SongId     string   //歌曲id
	Songname   string   //歌曲名
	Coverurl   string   //歌曲头像
	Username   string   //点歌用户名
	PlaylistID string   `gorm:"size:36"`               //当前歌曲属于哪个服务器的ID
	Playlist   Playlist `gorm:"foreignKey:PlaylistID"` //外键，指向Playlist的ID
}

func main() {

}

// 写入歌曲信息并创建播放列表
func CreateList(gid string, mid string, name string, user string, cover string) {

}
