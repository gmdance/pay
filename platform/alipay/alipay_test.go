package alipay

import (
	"encoding/json"
	"fmt"
	"testing"
)

func testError() (a []byte ,err error) {
	a, err = json.Marshal("aga")
	if err == nil {
		_, err = json.Marshal("{aga}")
	}
	return
}

func TestNewAlipay(t *testing.T) {
	_, err := json.Marshal("aga")
	if err == nil {
		fmt.Println(err)
		_, err := testError()
		fmt.Println(err)
	}
	fmt.Println(err)
}