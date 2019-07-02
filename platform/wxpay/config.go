package wxpay

const SignTypeMD5 = "MD5"
const SignTypeSHA256 = "HMAC-SHA256"

type Config struct {
	MchID           string `json:"mch_id"`
	Key             string `json:"key"`
	AppCertPem      string `json:"app_cert_pem"`
	AppKeyPem       string `json:"app_key_pem"`
	SignType        string `json:"sign_type"`
	PayNotifyURL    string `json:"pay_notify_url"`
	RefundNotifyURL string `json:"refund_notify_url"`
}
