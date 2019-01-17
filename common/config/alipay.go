package config

type AlipayConfig struct {
	AppID           string
	Partner         string
	SignType        string
	PublicKey       string
	PrivateKey      string
	PayNotifyURL    string
	RefundNotifyURL string
}
