package utils

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	url "net/url"
)

func HandleHttp(link string) string {
	//	将接收到的短链解析原始链接发送给到https处理
	//resp, err := http.Get("http://shorturl.8446666.sojson.com/parse/shorturl?url=" + link)

	return ""
}

func HandleHttps(link string) string {
	//判断链接是歌单还是单曲
	if link != "" {
		u, err := url.Parse(link)
		if err != nil {
			log.Error("识别链接类型有误！", err)
			return ""
		}
		if u == nil {
			return ""
		}
		path := u.Path
		if path == "playlist" {
			fmt.Println("当前链接中存在playlist")
		}
	} else {
		fmt.Println("链接为空 handlehttps")
	}

	return ""
}
