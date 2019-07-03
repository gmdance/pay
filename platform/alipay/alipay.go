package alipay

import (
	"bytes"
	"encoding/base64"
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
	conf           Config
	responseSuffix string
	errorResponse  string
	signNodeName   string
}

type Resp struct {
	Sign string `json:"sign"`
}



func NewAlipay(conf Config) (*Alipay) {
	return &Alipay{
		conf:           conf,
		responseSuffix: "_response",
		errorResponse:  "error_response",
		signNodeName:   "sign",
	}
}

func (alipay *Alipay) BuildQuery(method string, bizContent interface{}) (url.Values, error) {
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

func (alipay *Alipay) GetSignContent(params map[string]string) []byte {
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
	return buff.Bytes()
}

func (alipay *Alipay) SignParams(params map[string]string) (string, error) {
	return OpenSSLSign(alipay.GetSignContent(params), alipay.conf.SignType, alipay.conf.AppPrivateKey)
}

func (alipay *Alipay) ParseJsonResponse(responseContent, nodeName string, nodeIndex int) (string, string) {
	signDataStartIndex := nodeIndex + len(nodeName) + 2;
	signStartIndex := strings.LastIndex(responseContent, "\""+alipay.signNodeName+"\"")
	signDataEndIndex := signStartIndex - 1
	indexLen := signDataEndIndex - signDataStartIndex
	if indexLen < 0 {
		return "", ""
	}
	signData := responseContent[signDataStartIndex : signDataStartIndex+indexLen]
	sign := responseContent[signStartIndex:]
	signEndIndex := strings.LastIndex(sign, "\"}")
	sign = sign[8:signEndIndex]
	return signData, sign
}

func (alipay *Alipay) Request(method string, bizContent interface{}, resp interface{}) (string, error) {
	params, err := alipay.BuildQuery(method, bizContent)
	if err != nil {
		return "", err
	}
	raw, err := utils.HttpGet(MainHost + "?" + params.Encode())
	if err != nil {
		return "", err
	}
	data := string(raw)
	rootNodeName := strings.Replace(method, ".", "_", -1) + alipay.responseSuffix
	rootIndex := strings.Index(data, rootNodeName)
	errorIndex := strings.Index(data, alipay.errorResponse)
	content := ""
	sign := ""
	if rootIndex > 0 {
		content, sign = alipay.ParseJsonResponse(data, rootNodeName, rootIndex)
	} else {
		content, sign = alipay.ParseJsonResponse(data, alipay.errorResponse, errorIndex)
	}
	if alipay.conf.AlipayPublicKey != "" {
		signBytes, err := base64.StdEncoding.DecodeString(sign)
		err = OpenSSLVerify([]byte(content), signBytes, alipay.conf.SignType, alipay.conf.AlipayPublicKey)
		if err != nil {
			return data, err
		}
	}
	if resp != nil {
		err := json.Unmarshal([]byte(content), resp)
		if err != nil {
			return data, err
		}
	}
	return data, nil
}

