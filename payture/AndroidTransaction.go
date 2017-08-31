package payture

import (
	"net/http"
)

type AndroidTransaction struct {
	PayToken string
	Order    Payment
}

func (pay AndroidTransaction) Pay(merch Merchant) (*http.Response, error) {
	return pay.send("PAY", merch)
}

func (block AndroidTransaction) Block(merch Merchant) (*http.Response, error) {
	return block.send("BLOCK", merch)
}

func (trans AndroidTransaction) send(method string, merch Merchant) (*http.Response, error) {
	url := merch.Host + "/api/AndroidPay"
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["OrderId"] = []string{trans.Order.OrderId}
	params["Amount"] = []string{trans.Order.Amount}
	params["Method"] = []string{method}
	params["PayToken"] = []string{trans.PayToken}
	return sendRequest(url, params)
}
