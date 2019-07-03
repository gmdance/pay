package alipay

type Result struct {
	Code    string `json:"code"`
	Msg     string `json:"msg"`
	SubCode string `json:"sub_code"`
	SubMsg  string `json:"sub_msg"`
}

func (rs Result) IsSuccess() bool {
	return rs.Code == "100000"
}
