package alipay

import (
	"encoding/base64"
	"encoding/json"
)

type NotifyPayResp struct {
	AppId       string `json:"app_id"`
	SignType    string `json:"sign_type"`
	Sign        string `json:"sign"`
	TradeNo     string `json:"trade_no"`
	OutTradeNo  string `json:"out_trade_no"`
	TradeStatus string `json:"trade_status"`
	TotalAmount string `json:"total_amount"`
	RefundFee   string `json:"refund_fee"`
	Subject     string `json:"subject"`
	Body        string `json:"body"`
	BuyerId     string `json:"buyer_id"`
	SellerId    string `json:"seller_id"`
	NotifyId    string `json:"notify_id"`
	NotifyType  string `json:"notify_type"`
	NotifyTime  string `json:"notify_time"`
	Charset     string `json:"charset"`
	GmtCreate   string `json:"gmt_create"`
	GmtPayment  string `json:"gmt_payment"`
	GmtClose    string `json:"gmt_close"`
	Version     string `json:"version"`
}

func (alipay *Alipay) NotifyPay(params map[string]string) (*NotifyPayResp, error) {
	sign := params["sign"]
	signType := params["sign_type"]
	params["sign"] = ""
	params["sign_type"] = ""
	content := alipay.GetSignContent(params)
	signBytes, err := base64.StdEncoding.DecodeString(sign)
	if err != nil {
		return nil, err
	}
	err = OpenSSLVerify(content, signBytes, signType, alipay.conf.AlipayPublicKey)
	if err != nil {
		return nil, err
	}
	var res NotifyPayResp
	paramsBytes, _ := json.Marshal(params)
	_ = json.Unmarshal(paramsBytes, &res)
	return &res, nil
}

func (alipay *Alipay) NotifySuccess() string {
	return "success"
}

func (alipay *Alipay) NotifyFail() string {
	return "fail"
}
