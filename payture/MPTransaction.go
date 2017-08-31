package payture

import (
	"net/http"
)

type MPTransaction struct {
	Token string
	Order Payment
	CVC2  string
}

func (pay MPTransaction) Pay(merch Merchant) (*http.Response, error) {
	return pay.send("MPPay", merch)
}

func (block MPTransaction) Block(merch Merchant) (*http.Response, error) {
	return block.send("MPBlock", merch)
}

func (pay MPTransaction) send(cmd string, merch Merchant) (*http.Response, error) {
	url := merch.Host + "/api/" + cmd
	params := make(map[string][]string)
	params["Key"] = []string{merch.Key}
	params["OrderId"] = []string{pay.Order.OrderId}
	params["Amount"] = []string{pay.Order.Amount}
	params["CVC2"] = []string{pay.CVC2}
	params["Token"] = []string{pay.Token}
	return sendRequest(url, params)
}
