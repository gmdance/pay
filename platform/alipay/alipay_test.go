package alipay

import (
	"fmt"
	"strconv"
	"testing"
	"time"
)

var orderNo = strconv.FormatInt(time.Now().Unix(), 10)

func TestNewAlipay(t *testing.T) {
	conf := Config{
		AppID:           "2018041002529877",
		SignType:        SignTypeRSA2,
		AlipayPublicKey: "",
		AppPrivateKey: ``,
	}
	alipay := NewAlipay(conf)
	params := TradePreCreateParams{
		OutTradeNo:  orderNo,
		TotalAmount: "1",
		Subject:     "测试",
	}
	fmt.Println(alipay.TradePreCreate(params))
}
