package wxpay

import (
	"encoding/xml"
	"errors"
	"strconv"
)

//统一下单接口
type UnifiedOrderParams struct {
	AppID          string `xml:"app_id" json:"app_id"`                     //小程序或公众号appId 必填
	TradeType      string `xml:"trade_type" json:"trade_type"`             //交易类型 必填
	OutTradeNo     string `xml:"out_trade_no" json:"out_trade_no"`         //订单编号 必填
	TotalFee       int64  `xml:"total_fee" json:"total_fee"`               //金额 必填
	Body           string `xml:"body" json:"body"`                         //商品描述 必填
	SpbillCreateIp string `xml:"spbill_create_ip" json:"spbill_create_ip"` //终端ip 必填
	ProductID      string `xml:"product_id" json:"product_id"`             //商品id NATIVE必填
	OpenID         string `xml:"open_id" json:"open_id"`                   //OPENID JSAPI必填
	Detail         string `xml:"detail" json:"detail"`                     //商品详细 不必填
	FeeType        string `xml:"fee_type" json:"fee_type"`                 //货币种类 不必填 CNY
	DeviceInfo     string `xml:"device_info" json:"device_info"`           //自定义参数 不必填
	Attach         string `xml:"attach" json:"attach"`                     //附加数据 不必填
	TimeStart      string `xml:"time_start" json:"time_start"`             //交易起始时间 不必填
	TimeExpire     string `xml:"time_expire" json:"time_expire"`           //交易结束时间 不必填
	GoodsTag       string `xml:"goods_tag" json:"goods_tag"`               //订单优惠标记 不必填
	LimitPay       string `xml:"limit_pay" json:"limit_pay"`               //指定支付方式 不必填
	Receipt        string `xml:"receipt" json:"receipt"`                   //电子发票入口开放标识
	SceneInfo      struct {
		//场景信息
		Id       string `json:"id"`        //门店编号
		Name     string `json:"name"`      //门店名称
		AreaCode string `json:"area_code"` //门店行政区划码
		Address  string `json:"address"`   //门店详细地址
	} `xml:"scene_info" json:"scene_info"`
}

type UnifiedOrderResp struct {
	WxpayError
	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

func (wxpay *Wxpay) UnifiedOrder(order UnifiedOrderParams) (*UnifiedOrderResp, []byte, error) {
	if order.AppID == "" {
		return nil, nil, errors.New("appId未填写")
	}
	if order.Body == "" {
		return nil, nil, errors.New("body未填写")
	}
	if order.OutTradeNo == "" {
		return nil, nil, errors.New("orderNo未填写")
	}
	if order.TotalFee == 0 {
		return nil, nil, errors.New("amount未填写")
	}
	if order.SpbillCreateIp == "" {
		return nil, nil, errors.New("clientIp未填写")
	}
	if order.TradeType == "" {
		return nil, nil, errors.New("tradeType未填写")
	}
	if order.TradeType == WxpayTradeTypeNative {
		if order.ProductID == "" {
			return nil, nil, errors.New("productId未填写")
		}
	} else if order.TradeType == WxpayTradeTypeJsapi {
		if order.OpenID == "" {
			return nil, nil, errors.New("openId未填写")
		}
	}

	data := map[string]string{
		"appid":            order.AppID,
		"body":             order.Body,
		"detail":           order.Detail,
		"out_trade_no":     order.OutTradeNo,
		"total_fee":        strconv.FormatInt(order.TotalFee, 10),
		"spbill_create_ip": order.SpbillCreateIp,
		"notify_url":       wxpay.conf.PayNotifyURL,
		"trade_type":       order.TradeType,
		"product_id":       order.ProductID,
		"openid":           order.OpenID,
	}
	var response UnifiedOrderResp
	_, raw, err := wxpay.Request(UriPathUnifiedOrder, data)
	xml.Unmarshal(raw, &response)
	return &response, raw, err
}
