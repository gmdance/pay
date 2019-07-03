package wxpay

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
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

	MainHost            = "https://api.mch.weixin.qq.com"
	UriPathUnifiedOrder = "/pay/unifiedorder"
	UriPathRefund       = "/secapi/pay/refund"
	UriPathOrderQuery   = "/pay/orderquery"

	WxpayTradeTypeNative = "NATIVE"
	WxpayTradeTypeJsapi  = "JSAPI"
	WxpayTradeTypeApp    = "APP"
)

type Wxpay struct {
	conf Config
}

type WxpayResp struct {
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

func NewWxpay(conf Config) (*Wxpay) {
	if conf.SignType == "" {
		conf.SignType = "MD5"
	}
	return &Wxpay{
		conf: conf,
	}
}

func (wxpay *Wxpay) SignParams(data map[string]string) string {
	var keys []string
	for k := range data {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var buff bytes.Buffer
	for _, k := range keys {
		value := data[k]
		value = strings.TrimSpace(value)
		if value == "" || k == "sign" {
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
	buff.WriteString(wxpay.conf.Key)
	sign := ""
	if wxpay.conf.SignType == SignTypeMD5 {
		m := md5.New()
		m.Write(buff.Bytes())
		sign = hex.EncodeToString(m.Sum(nil))
		sign = strings.ToUpper(sign)
	} else if wxpay.conf.SignType == SignTypeSHA256 {
		h := hmac.New(sha256.New, []byte(wxpay.conf.Key))
		_, _ = io.WriteString(h, buff.String())
		sign = hex.EncodeToString(h.Sum(nil))
		sign = strings.ToUpper(sign)
	}
	return sign
}

func (wxpay *Wxpay) Request(api string, params map[string]string, resp interface{}) (data string, e error) {
	data = ""
	if wxpay.conf.MchID == "" {
		return data, errors.New("mchId未配置")
	}
	if wxpay.conf.SignType == "" {
		return data, errors.New("signType未配置")
	}
	if wxpay.conf.Key == "" {
		return data, errors.New("wxKey未配置")
	}
	apiURL := MainHost + path.Join("/", api)
	params["mch_id"] = wxpay.conf.MchID
	params["sign_type"] = wxpay.conf.SignType
	params["nonce_str"] = strconv.Itoa(rand.New(rand.NewSource(time.Now().Unix())).Int())
	sign := wxpay.SignParams(params)
	params["sign"] = sign
	raw, err := xml.Marshal(utils.Xml(params))
	if err != nil {
		return data, err
	}
	body, err := utils.HttpPost(apiURL, "application/xml", raw)
	if err != nil {
		return data, err
	}
	data = string(body)
	resultMap := make(map[string]string)
	err = xml.Unmarshal(body, (*utils.Xml)(&resultMap))
	if err != nil {
		return data, err
	}
	if resultMap["return_code"] != WxpaySuccess {
		return data, errors.New("微信通讯失败:" + resultMap["return_msg"])
	}
	checkSign := wxpay.SignParams(resultMap)
	if checkSign != resultMap["sign"] {
		return data, errors.New("微信返回签名失败")
	}
	if resultMap["result_code"] != WxpaySuccess {
		return data, errors.New(fmt.Sprintf("微信业务失败:%s(%s)", resultMap["err_code_des"], resultMap["err_code"]))
	}
	if resp != nil {
		e = xml.Unmarshal(raw, &resp)
	}
	return
}
