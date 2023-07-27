package conf

const (
	Token   = "1/MjEyMzY=/51EiCA37ybJqwbhVlY3NmQ=="
	BaseUrl = "https://www.kookapp.cn/api"
	NetEasy = "http://47.96.25.105:3000" //网易云的api
	// HTTPServerIp HTTPServerPort  VerifyToken EncryptKey : WEBHOOK相关, 如果不是WEBHOOK,可不填
	// HTTPServerIp 侦听的ip地址
	HTTPServerIp = "0.0.0.0"
	// HTTPServerPort 侦听的端口
	HTTPServerPort = "8080"
	// VerifyToken 若不需要verify, 可不填。填上可以防止别人填你的地址，导致你的消息混乱。
	VerifyToken = ""
	// EncryptKey 若不需要加密，可不填。如果有encryptKey会更安全一点。
	EncryptKey = ""
	Version    = "v0.0.5"
	//数据库相关配置
	Databese = "kookbot"
	Username = "kookbot"
	Host     = "mysql.sqlpub.com:3306"
	Password = "18e95fb4e7f1cba7"
	//数据库相关配置结束
)
