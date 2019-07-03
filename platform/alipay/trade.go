package alipay

type TradePreCreateParams struct {
	OutTradeNo           string `json:"out_trade_no"`            //订单号 必填
	TotalAmount          string `json:"total_amount"`            //订单金额 必填
	Subject              string `json:"subject"`                 //订单标题 必填
	SellerId             string `json:"seller_id"`               //支付宝用户ID
	DiscountableAmount   string `json:"discountable_amount"`     //可打折金额
	Body                 string `json:"body"`                    //对商品的描述
	ProductCode          string `json:"product_code"`            //销售产品码
	OperatorId           string `json:"operator_id"`             //商户操作员编码
	StoreId              string `json:"store_id"`                //商户门店编码
	DisablePayChannels   string `json:"disable_pay_channels"`    //禁止渠道
	EnablePayChannels    string `json:"enable_pay_channels"`     //可用渠道
	TerminalId           string `json:"terminal_id"`             //终端id
	TimeoutExpress       string `json:"timeout_express"`         //该笔订单允许的最晚付款时间，逾期将关闭交易
	MerchantOrderNo      string `json:"merchant_order_no"`       //商户原始订单号
	QrCodeTimeoutExpress string `json:"qr_code_timeout_express"` //该笔订单最晚付款时间
}

type TradePreCreateResult struct {
	Result
	OutTradeNo string `json:"out_trade_no"` //商户订单号
	QrCode     string `json:"qr_code"`      //支付二维码
}

func (alipay *Alipay) TradePreCreate(bizContent TradePreCreateParams) (*TradePreCreateResult, string, error) {
	var result TradePreCreateResult
	data, err := alipay.Request(MethodTradePreCreate, bizContent, &result)
	return &result, data, err
}
