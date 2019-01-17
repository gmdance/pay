package api

import (
	"github.com/gmdance/pay/common/config"
	"strconv"
	"testing"
	"time"
)

func TestNewWechatApiUnifiedOrder(t *testing.T) {
	conf := config.WechatConfig{
		MchID:        "1900009851",
		Key:          "8934e7d15453e97507ef794cf7b0519d",
		SignType:     "MD5",
		PayNotifyURL: "https://notify.yidu.ai/notify/wechat",
	}
	wechatApi := NewWechatApi(conf)
	order := WechatApiOrder{
		AppID:     "wx426b3015555a46be",
		TradeType: WechatTradeTypeNative,
		OrderNo:   strconv.FormatInt(time.Now().Unix(), 10),
		Amount:    1,
		Currency:  "CNY",
		Body:      "test",
		ClientIP:  "127.0.0.1",
		ProductID: "1",
		OpenID:    "",
		Detail:    "detail",
	}
	response, _, err := wechatApi.UnifiedOrder(order)
	t.Log(response)
	if err != nil {
		t.Fail()
	}
}
