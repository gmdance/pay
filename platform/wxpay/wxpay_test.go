package wxpay

import (
	"strconv"
	"testing"
	"time"
)

var (
	conf = Config{
		MchID:        "1900009851",
		Key:          "8934e7d15453e97507ef794cf7b0519d",
		SignType:     "MD5",
		PayNotifyURL: "https://notify.yidu.ai/notify/wechat",
	}
	orderNo = strconv.FormatInt(time.Now().Unix(), 10)
)

func TestNewWechatApiUnifiedOrder(t *testing.T) {
	wechatApi := NewWxpay(conf)
	order := UnifiedOrderParams{
		AppID:          "wx426b3015555a46be",
		TradeType:      WxpayTradeTypeNative,
		OutTradeNo:     orderNo,
		TotalFee:       1,
		FeeType:        "CNY",
		Body:           "test",
		SpbillCreateIp: "127.0.0.1",
		ProductID:      "1",
		OpenID:         "",
		Detail:         "detail",
	}
	response, raw, err := wechatApi.UnifiedOrder(order)
	t.Log(string(raw))
	t.Log(response)
	if err != nil {
		t.Fail()
	}
}

func TestWechatApi_OrderQuery(t *testing.T) {
	wechatApi := NewWxpay(conf)
	response, raw, err := wechatApi.OrderQuery("wx426b3015555a46be", orderNo, "")
	t.Log(string(raw))
	t.Log(response)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
}
