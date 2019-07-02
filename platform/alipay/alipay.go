package alipay

import (
	"bytes"
	"encoding/json"
	"github.com/gmdance/pay/utils"
	"net/url"
	"sort"
	"strings"
	"time"
)

const (
	TradeTypeWeb = "web"
	TradeTypeWap = "wap"
	TradeTypeQR  = "qr"
	TradeTypeApp = "app"

	SignTypeRSA  = "RSA"
	SignTypeRSA2 = "RSA2"

	MainHost             = "https://openapi.alipay.com/gateway.do"
	MethodTradePreCreate = "alipay.trade.precreate"
	MethodTradeQuery     = "alipay.trade.query"
)

type Alipay struct {
	conf Config
}

type AlipayRequest struct {
	AppID        string `json:"app_id"`
	Method       string `json:"method"`
	Charset      string `json:"charset"`
	SignType     string `json:"sign_type"`
	Sign         string `json:"sign"`
	Timestamp    string `json:"timestamp"`
	Version      string `json:"version"`
	NotifyUrl    string `json:"notify_url"`
	AppAuthToken string `json:"app_auth_token"`
	BizContent   string `json:"biz_content"`
}

type AlipayError struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
	Sign    string `json:"sign"`
}

func NewAlipay(conf Config) (*Alipay) {
	return &Alipay{
		conf: conf,
	}
}

func (alipay *Alipay) BuildParams(method string, bizContent interface{}) (url.Values, error) {
	conf := alipay.conf
	bizContentData, err := json.Marshal(bizContent)
	if err != nil {
		return nil, err
	}
	params := map[string]string{
		"app_id":      conf.AppID,
		"method":      method,
		"format":      "JSON",
		"charset":     "utf-8",
		"sign_type":   conf.SignType,
		"timestamp":   time.Now().Format("2006-01-02 03:04:05"),
		"version":     "1.0",
		"notify_url":  conf.PayNotifyURL,
		"biz_content": string(bizContentData),
	}
	sign, err := alipay.SignParams(params)
	if err != nil {
		return nil, err
	}
	params["sign"] = sign
	values := url.Values{}
	for k, v := range params {
		values.Set(k, v)
	}
	return values, nil
}

func (alipay *Alipay) SignParams(params map[string]string) (string, error) {
	keys := make([]string, 0)
	for key, _ := range params {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buff bytes.Buffer
	for _, k := range keys {
		value := params[k]
		value = strings.TrimSpace(value)
		if params[k] == "" {
			continue
		}
		if buff.Len() != 0 {
			buff.WriteString("&")
		}
		buff.WriteString(k)
		buff.WriteString("=")
		buff.WriteString(value)
	}
	return Sign(buff.Bytes(), alipay.conf.SignType, alipay.conf.PrivateKey)
}

func (alipay *Alipay) Request(method string, bizContent interface{}, resp interface{}) (string, error) {
	params, err := alipay.BuildParams(method, bizContent)
	if err != nil {
		return "", err
	}
	raw, err := utils.Get(MainHost + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	return string(raw), nil
}
