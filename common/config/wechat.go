package config

const WxSignTypeMD5 = "MD5"
const WxSignTypeSHA256 = "HMAC-SHA256"

type WechatConfig struct {
	MchID           string
	Key             string
	AppCertPem      string
	AppKeyPem       string
	SignType        string
	PayNotifyURL    string
	RefundNotifyURL string
}
