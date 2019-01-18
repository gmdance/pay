package api

import "github.com/gmdance/pay/config"

type AlipayApi struct {
	conf config.AlipayConfig
}

func NewAlipayApi(conf config.AlipayConfig) (*AlipayApi) {
	return &AlipayApi{
		conf: conf,
	}
}
