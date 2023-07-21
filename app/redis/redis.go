package redisDB

import "github.com/redis/go-redis/v9"

var Rdb *redis.Client

func init() {
	//	初始化redis配置
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

type PlayList struct {
	List_id string `json:"list_id"`
}

func CreateList(p *PlayList) {
	//	当频道没有歌单时创建歌单

}
