package api

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/gmdance/pay/config"
	"github.com/gmdance/pay/utils"
	"io"
	"math/rand"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	WxpaySuccess              = "SUCCESS"
	WxpayFail                 = "FAIL"
	WxpayTradeStateSuccess    = "SUCCESS"
	WxpayTradeStateRefund     = "REFUND"
	WxpayTradeStateNopay      = "NOTPAY"
	WxpayTradeStateClosed     = "CLOSED"
	WxpayTradeStateRevoked    = "REVOKED"
	WxpayTradeStateUserPaying = "USERPAYING"
	WxpayTradeStatePayError   = "PAYERROR"

	wxpayHost             = "https://api.mch.weixin.qq.com"
	wxpayPathUnifiedOrder = "/pay/unifiedorder"
	wxpayPathRefund       = "/secapi/pay/refund"
	wxpayPathOrderQuery   = "/pay/orderquery"

	WxpayTradeTypeNative = "NATIVE"
)

type WxpayApi struct {
	conf config.WxpayConfig
}

type WxpayError struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	ResultCode string `xml:"result_code"`
	ResultMsg  string `xml:"result_msg"`
	ErrCode    string `xml:"err_code"`
	ErrCodeDes string `xml:"err_code_des"`
	AppID      string `xml:"appid"`
	MchID      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Sign       string `xml:"sign"`
	SignType   string `xml:"sign_type"`
	DeviceInfo string `xml:"device_info"`
}

func NewWxpayApi(conf config.WxpayConfig) (*WxpayApi) {
	if conf.SignType == "" {
		conf.SignType = "MD5"
	}
	return &WxpayApi{
		conf: conf,
	}
}

func (wa *WxpayApi) Sign(data map[string]string) string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buff bytes.Buffer
	for _, k := range keys {
		value := data[k]
		if data[k] == "" || k == "sign" {
			continue
		}
		if buff.Len() != 0 {
			buff.WriteString("&")
		}
		buff.WriteString(k)
		buff.WriteString("=")
		buff.WriteString(value)
	}
	buff.WriteString("&key=")
	buff.WriteString(wa.conf.Key)
	sign := ""
	if wa.conf.SignType == config.WxSignTypeMD5 {
		m := md5.New()
		m.Write(buff.Bytes())
		sign = hex.EncodeToString(m.Sum(nil))
		sign = strings.ToUpper(sign)
	} else if wa.conf.SignType == config.WxSignTypeSHA256 {
		h := hmac.New(sha256.New, []byte(wa.conf.Key))
		io.WriteString(h, buff.String())
		sign = hex.EncodeToString(h.Sum(nil))
		sign = strings.ToUpper(sign)
	}
	return sign
}

func (wa *WxpayApi) Post(api string, data map[string]string) (map[string]string, []byte, error) {
	apiURL := wxpayHost + path.Join("/", api)
	data["mch_id"] = wa.conf.MchID
	data["sign_type"] = wa.conf.SignType
	data["nonce_str"] = strconv.Itoa(rand.New(rand.NewSource(time.Now().Unix())).Int())
	sign := wa.Sign(data)
	data["sign"] = sign
	raw, err := xml.Marshal(utils.Xml(data))
	if err != nil {
		return nil, nil, err
	}
	body, err := utils.Post(apiURL, "application/xml", raw)
	if err != nil {
		return nil, nil, err
	}
	resultMap := make(map[string]string)
	err = xml.Unmarshal(body, (*utils.Xml)(&resultMap))
	if err != nil {
		return nil, body, err
	}
	if resultMap["return_code"] != WxpaySuccess {
		return resultMap, body, errors.New("微信通讯失败:" + resultMap["return_msg"])
	}
	checkSign := wa.Sign(resultMap)
	if checkSign != resultMap["sign"] {
		return resultMap, body, errors.New("微信返回签名失败")
	}
	if resultMap["result_code"] != WxpaySuccess {
		return resultMap, body, errors.New(fmt.Sprintf("微信业务失败:%s(%s)", resultMap["err_code_des"], resultMap["err_code"]))
	}
	return resultMap, body, nil
}

//统一下单接口
type WxpayUnifiedOrderRequest struct {
	AppID     string //小程序或公众号appId 必填
	TradeType string //交易类型 必填
	OrderNo   string //订单编号 必填
	Amount    int64  //金额 必填
	Currency  string //货币种类 不必填
	Body      string //商品描述 必填
	ClientIP  string //终端ip 必填
	ProductID string //商品id NATIVE必填
	OpenID    string //OPENID JSAPI必填
	Detail    string //商品详细 不必填
}

type WxpayUnifiedOrderResponse struct {
	WxpayError
	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

func (wa *WxpayApi) UnifiedOrder(order WxpayUnifiedOrderRequest) (WxpayUnifiedOrderResponse, []byte, error) {
	data := map[string]string{
		"appid":            order.AppID,
		"body":             order.Body,
		"detail":           order.Detail,
		"out_trade_no":     order.OrderNo,
		"total_fee":        strconv.FormatInt(order.Amount, 10),
		"spbill_create_ip": order.ClientIP,
		"notify_url":       wa.conf.PayNotifyURL,
		"trade_type":       order.TradeType,
		"product_id":       order.ProductID,
		"openid":           order.OpenID,
	}
	var response WxpayUnifiedOrderResponse
	_, raw, err := wa.Post(wxpayPathUnifiedOrder, data)
	xml.Unmarshal(raw, &response)
	return response, raw, err
}

//查询订单接口
type WxpayOrderQueryResponse struct {
	WxpayError
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

func (wa *WxpayApi) OrderQuery(appID, orderNo, transactionId string) (WxpayOrderQueryResponse, []byte, error) {
	data := map[string]string{
		"appid":          appID,
		"out_trade_no":   orderNo,
		"transaction_id": transactionId,
	}
	_, raw, err := wa.Post(wxpayPathOrderQuery, data)
	var response WxpayOrderQueryResponse
	xml.Unmarshal(raw, &response)
	return response, raw, err
}

//退款接口
type WxpayRefundRequest struct {
	AppID        string
	OrderNo      string
	RefundNo     string
	OrderAmount  int64
	RefundAmount int64
	Currency     string
	RefundDesc   string
}

func (wa *WxpayApi) Refund(refund WxpayRefundRequest) (map[string]string, []byte, error) {
	data := map[string]string{
		"appid":         refund.AppID,
		"out_trade_no":  refund.OrderNo,
		"out_refund_no": refund.RefundNo,
		"total_fee":     strconv.FormatInt(refund.OrderAmount, 10),
		"refund_fee":    strconv.FormatInt(refund.RefundAmount, 10),
		"notify_url":    wa.conf.RefundNotifyURL,
		"refund_desc":   refund.RefundDesc,
	}
	return wa.Post(wxpayPathRefund, data)
}

//支付回调校验
type WechatApiPayNotifyData struct {
	WxpayError
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

func (wa *WxpayApi) PayNotify(raw string) (*WechatApiPayNotifyData, error) {
	rawBytes := []byte(raw)
	data := make(map[string]string)
	err := xml.Unmarshal(rawBytes, (*utils.Xml)(&data))
	if err != nil {
		return nil, err
	}
	sign := wa.Sign(data)
	if sign != data["sign"] {
		return nil, errors.New("微信支付回调签名失败")
	}
	var notifyData WechatApiPayNotifyData
	err = xml.Unmarshal(rawBytes, &notifyData)
	return &notifyData, err
}

//回调成功返回
func (wa *WxpayApi) NotifySuccess() string {
	return "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
}
