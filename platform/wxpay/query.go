package wxpay

import (
	"encoding/xml"
	"errors"
)

//查询订单接口
type OrderQueryResp struct {
	WxpayResp
	OpenID             string `xml:"open_id"`
	IsSubscribe        string `xml:"is_subscribe"`
	TradeType          string `xml:"trade_type"`
	TradeState         string `xml:"trade_state"`
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

func (wxpay *Wxpay) OrderQuery(appID, orderNo, transactionId string) (*OrderQueryResp, []byte, error) {
	if orderNo == "" && transactionId == "" {
		return nil, nil, errors.New("orderNo和transactionId必须填写一项")
	}
	data := map[string]string{
		"appid":          appID,
		"out_trade_no":   orderNo,
		"transaction_id": transactionId,
	}
	_, raw, err := wxpay.Request(UriPathOrderQuery, data)
	var response OrderQueryResp
	xml.Unmarshal(raw, &response)
	return &response, raw, err
}