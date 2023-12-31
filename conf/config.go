package conf

var Token = ""

const (
	BaseUrl    = "https://www.kookapp.cn/api"
	NetEasy    = "http://47.96.25.105:3000" //网易云的api
	OnlineUUID = ""                         //BotMarket验证在线uuid,可选
	// HTTPServerIp HTTPServerPort  VerifyToken EncryptKey : WEBHOOK相关, 如果不是WEBHOOK,可不填
	// HTTPServerIp 侦听的ip地址
	HTTPServerIp = "0.0.0.0"
	// HTTPServerPort 侦听的端口
	HTTPServerPort = "8080"
	// VerifyToken 若不需要verify, 可不填。填上可以防止别人填你的地址，导致你的消息混乱。
	VerifyToken = ""
	// EncryptKey 若不需要加密，可不填。如果有encryptKey会更安全一点。
	EncryptKey = ""
	Version    = "v0.1.3"
	//数据库相关配置
	Databese = "test"
	Username = "x1a2h1"
	Host     = "192.168.110.69:3306"
	Password = "Qkxz1216"
	//数据库相关配置结束
)
const envToken = ""

//func init() {
//	Token = os.Getenv("TOKEN")
//}
