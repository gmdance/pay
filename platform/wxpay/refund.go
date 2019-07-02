package wxpay

import "strconv"

//退款接口
type RefundParams struct {
	AppID        string
	OrderNo      string
	RefundNo     string
	OrderAmount  int64
	RefundAmount int64
	Currency     string
	RefundDesc   string
}

func (wxpay *Wxpay) Refund(refund RefundParams) (map[string]string, []byte, error) {
	data := map[string]string{
		"appid":         refund.AppID,
		"out_trade_no":  refund.OrderNo,
		"out_refund_no": refund.RefundNo,
		"total_fee":     strconv.FormatInt(refund.OrderAmount, 10),
		"refund_fee":    strconv.FormatInt(refund.RefundAmount, 10),
		"notify_url":    wxpay.conf.RefundNotifyURL,
		"refund_desc":   refund.RefundDesc,
	}
	return wxpay.Request(UriPathRefund, data)
}
