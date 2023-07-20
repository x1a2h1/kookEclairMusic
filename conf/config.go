package conf

const (
	Token   = "1/MTI4NTk=/aYonvRLa4CMd4yuvVh19hg=="
	BaseUrl = "https://www.kookapp.cn/api"
	NetEasy = "http://192.168.110.69:3000" //网易云的api
	// HTTPServerIp HTTPServerPort  VerifyToken EncryptKey : WEBHOOK相关, 如果不是WEBHOOK,可不填
	// HTTPServerIp 侦听的ip地址
	HTTPServerIp = "0.0.0.0"
	// HTTPServerPort 侦听的端口
	HTTPServerPort = "8080"
	// VerifyToken 若不需要verify, 可不填。填上可以防止别人填你的地址，导致你的消息混乱。
	VerifyToken = ""
	// EncryptKey 若不需要加密，可不填。如果有encryptKey会更安全一点。
	EncryptKey = ""
	Version    = "v0.0.1"
)