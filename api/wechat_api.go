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
	wechatSuccess              = "SUCCESS"
	wechatFail                 = "FAIL"
	WechatTradeStateSuccess    = "SUCCESS"
	WechatTradeStateRefund     = "REFUND"
	WechatTradeStateNopay      = "NOTPAY"
	WechatTradeStateClosed     = "CLOSED"
	WechatTradeStateRevoked    = "REVOKED"
	WechatTradeStateUserPaying = "USERPAYING"
	WechatTradeStatePayError   = "PAYERROR"

	wechatHost             = "https://api.mch.weixin.qq.com"
	wechatPathUnifiedOrder = "/pay/unifiedorder"
	wechatPathRefund       = "/secapi/pay/refund"
	wechatPathOrderQuery   = "/pay/orderquery"

	WechatTradeTypeNative = "NATIVE"
)

type wechatApi struct {
	conf config.WechatConfig
}

type wechatError struct {
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

func NewWechatApi(conf config.WechatConfig) (*wechatApi) {
	if conf.SignType == "" {
		conf.SignType = "MD5"
	}
	return &wechatApi{
		conf: conf,
	}
}

func (wa *wechatApi) sign(data map[string]string) string {
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

func (wa *wechatApi) Post(api string, data map[string]string) (map[string]string, []byte, error) {
	apiURL := wechatHost + path.Join("/", api)
	data["mch_id"] = wa.conf.MchID
	data["sign_type"] = wa.conf.SignType
	data["nonce_str"] = strconv.Itoa(rand.New(rand.NewSource(time.Now().Unix())).Int())
	sign := wa.sign(data)
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
	if resultMap["return_code"] != wechatSuccess {
		return resultMap, body, errors.New("微信通讯失败:" + resultMap["return_msg"])
	}
	checkSign := wa.sign(resultMap)
	if checkSign != resultMap["sign"] {
		return resultMap, body, errors.New("微信返回签名失败")
	}
	if resultMap["result_code"] != wechatSuccess {
		return resultMap, body, errors.New(fmt.Sprintf("微信业务失败:%s(%s)", resultMap["err_code_des"], resultMap["err_code"]))
	}
	return resultMap, body, nil
}

//统一下单接口
type WechatApiOrder struct {
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

type WechatUnifiedOrderResponse struct {
	wechatError
	TradeType string `xml:"trade_type"`
	PrepayID  string `xml:"prepay_id"`
	CodeURL   string `xml:"code_url"`
}

func (wa *wechatApi) UnifiedOrder(order WechatApiOrder) (WechatUnifiedOrderResponse, []byte, error) {
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
	var response WechatUnifiedOrderResponse
	_, raw, err := wa.Post(wechatPathUnifiedOrder, data)
	xml.Unmarshal(raw, &response)
	return response, raw, err
}

//查询订单接口
type WechatOrderQueryResponse struct {
	wechatError
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

func (wa *wechatApi) OrderQuery(appID, orderNo, transactionId string) (WechatOrderQueryResponse, []byte, error) {
	data := map[string]string{
		"appid":          appID,
		"out_trade_no":   orderNo,
		"transaction_id": transactionId,
	}
	_, raw, err := wa.Post(wechatPathOrderQuery, data)
	var response WechatOrderQueryResponse
	xml.Unmarshal(raw, &response)
	return response, raw, err
}

//退款接口
type WechatApiRefund struct {
	AppID        string
	OrderNo      string
	RefundNo     string
	OrderAmount  int64
	RefundAmount int64
	Currency     string
	RefundDesc   string
}

func (wa *wechatApi) Refund(refund WechatApiRefund) {
	data := map[string]string{
		"appid":         refund.AppID,
		"out_trade_no":  refund.OrderNo,
		"out_refund_no": refund.RefundNo,
		"total_fee":     strconv.FormatInt(refund.OrderAmount, 10),
		"refund_fee":    strconv.FormatInt(refund.RefundAmount, 10),
		"notify_url":    wa.conf.RefundNotifyURL,
		"refund_desc":   refund.RefundDesc,
	}
	wa.Post(wechatPathRefund, data)
}

//支付回调校验
type WechatApiPayNotifyData struct {
	wechatError
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

func (wa *wechatApi) PayNotify(raw string) (*WechatApiPayNotifyData, error) {
	rawBytes := []byte(raw)
	data := make(map[string]string)
	err := xml.Unmarshal(rawBytes, (*utils.Xml)(&data))
	if err != nil {
		return nil, err
	}
	sign := wa.sign(data)
	if sign != data["sign"] {
		return nil, errors.New("微信支付回调签名失败")
	}
	var notifyData WechatApiPayNotifyData
	err = xml.Unmarshal(rawBytes, &notifyData)
	return &notifyData, err
}

func (wa *wechatApi) NotifySuccess() string {
	return "<xml><return_code><![CDATA[SUCCESS]]></return_code><return_msg><![CDATA[OK]]></return_msg></xml>"
}
