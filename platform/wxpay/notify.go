package wxpay

import (
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gmdance/pay/utils"
)

//支付回调校验
type NotifyPayResp struct {
	WxpayResp
	OpenID             string `xml:"openid"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	BankType           string `xml:"bank_type"`
	TotalFee           string `xml:"total_fee"`
	SettlementTotalFee string `xml:"settlement_total_fee"`
	FeeType            string `xml:"fee_type"`
	CashFee            string `xml:"cash_fee"`
	CashFeeType        string `xml:"cash_fee_type"`
	TransactionID      string `xml:"transaction_id"`
	OutTradeNo         string `xml:"out_trade_no"`
	TimeEnd            string `xml:"time_end"`
}

//支付回调校验
func (wxpay *Wxpay) NotifyPay(raw string) (*NotifyPayResp, error) {
	rawBytes := []byte(raw)
	data := make(map[string]string)
	err := xml.Unmarshal(rawBytes, (*utils.Xml)(&data))
	if err != nil {
		return nil, err
	}
	sign := wxpay.SignParams(data)
	if sign != data["sign"] {
		return nil, errors.New("微信支付回调签名失败")
	}
	var notifyData NotifyPayResp
	err = xml.Unmarshal(rawBytes, &notifyData)
	return &notifyData, err
}

//回调成功返回
func (wxpay *Wxpay) NotifySuccess() string {
	return "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
}

//回调失败返回
func (wxpay *Wxpay) NotifyFail(message string) string {
	return fmt.Sprintf("<xml><return_code><![CDATA[FAIL]]></return_code><return_msg><![CDATA[%s]]></return_msg></xml>", message)
}
