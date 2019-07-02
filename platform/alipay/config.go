package alipay

type Config struct {
	AppID           string `json:"app_id"`
	Partner         string `json:"partner"`
	SignType        string `json:"sign_type"`
	PublicKey       string `json:"public_key"`
	PrivateKey      string `json:"private_key"`
	PayNotifyURL    string `json:"pay_notify_url"`
	RefundNotifyURL string `json:"refund_notify_url"`
}

