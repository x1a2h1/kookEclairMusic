package netease

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func HandleLink(link string) (string, string, error) {
	//	判断当前链接是否是原始链接

	shortURL := ""
	if strings.HasPrefix(link, "http://163cn.tv/") {
		client := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return http.ErrUseLastResponse // 不自动处理重定向，这样我们可以获取到重定向的URL
			},
		}
		fmt.Println("当前是短链接")
		resp, err := client.Get(link)
		if err != nil {
			return "", "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 300 && resp.StatusCode < 400 {
			loc, err := resp.Location()
			if err != nil {
				return "", "", err
			}
			shortURL = loc.String()
		}
	} else {
		shortURL = link
	}
	fmt.Println("当前链接为：", shortURL)
	//	将链接进行处理判断song？orplaylist？ 获取他们的id
	parsedURL, err := url.Parse(shortURL)
	if err != nil {
		return "", "", err
	}
	path := parsedURL.Path
	params := ""
	switch {
	case path == "/playlist":
		params = "playlist"
	case path == "/song":
		params = "song"
	default:
		return "", "", err
	}
	id := parsedURL.Query().Get("id")
	if id != "" {
		return params, id, err
	} else {
		return "", "", err
	}
}
