package payture

import (
	"net/http"
)

type AppleTransaction struct {
	PayToken string
	OrderId  string
}

func (pay AppleTransaction) Pay(merch Merchant) (*http.Response, error) {
	return pay.send("PAY", merch)
}

func (block AppleTransaction) Block(merch Merchant) (*http.Response, error) {
	return block.send("BLOCK", merch)
}

func (trans AppleTransaction) send(method string, merch Merchant) (*http.Response, error) {
	url := merch.Host + "/api/ApplePay"
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["OrderId"] = []string{trans.OrderId}
	params["Method"] = []string{method}
	params["PayToken"] = []string{trans.PayToken}
	return sendRequest(url, params)
}
