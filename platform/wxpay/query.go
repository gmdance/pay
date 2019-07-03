package wxpay

import (
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

func (wxpay *Wxpay) OrderQuery(appID, orderNo, transactionId string) (*OrderQueryResp, string, error) {
	if orderNo == "" && transactionId == "" {
		return nil, "", errors.New("orderNo和transactionId必须填写一项")
	}
	params := map[string]string{
		"appid":          appID,
		"out_trade_no":   orderNo,
		"transaction_id": transactionId,
	}
	var response OrderQueryResp
	data, err := wxpay.Request(UriPathOrderQuery, params, &response)
	return &response, data, err
}