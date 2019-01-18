package api

import (
	"github.com/gmdance/pay/config"
	"strconv"
	"testing"
	"time"
)

var (
	conf = config.WxpayConfig{
		MchID:        "1900009851",
		Key:          "8934e7d15453e97507ef794cf7b0519d",
		SignType:     "MD5",
		PayNotifyURL: "https://notify.yidu.ai/notify/wechat",
	}
	orderNo = strconv.FormatInt(time.Now().Unix(), 10)
)

func TestNewWechatApiUnifiedOrder(t *testing.T) {
	wechatApi := NewWxpayApi(conf)
	order := WxpayUnifiedOrderRequest{
		AppID:     "wx426b3015555a46be",
		TradeType: WxpayTradeTypeNative,
		OrderNo:   orderNo,
		Amount:    1,
		Currency:  "CNY",
		Body:      "test",
		ClientIP:  "127.0.0.1",
		ProductID: "1",
		OpenID:    "",
		Detail:    "detail",
	}
	response, raw, err := wechatApi.UnifiedOrder(order)
	t.Log(string(raw))
	t.Log(response)
	if err != nil {
		t.Fail()
	}
}

func TestWechatApi_OrderQuery(t *testing.T) {
	wechatApi := NewWxpayApi(conf)
	response, raw, err := wechatApi.OrderQuery("wx426b3015555a46be", orderNo, "")
	t.Log(string(raw))
	t.Log(response)
	if err != nil {
		t.Log(err.Error())
		t.Fail()
	}
}
