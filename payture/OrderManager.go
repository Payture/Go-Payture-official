package payture

import "fmt"

type OrderManager struct {
	apiType  string
	merchant Merchant
}

func GetOrderManager(api string, merch Merchant) (ordM OrderManager) {
	ordM.apiType = api
	ordM.merchant = merch
	return
}

func (this OrderManager) getReqUrl(cmd string) string {
	return this.merchant.Host + "/" + this.apiType + "/" + cmd
}

func (this OrderManager) Unblock(order Payment) (orderResp OrderResponse, err error) {
	key := "Key"
	if this.apiType == "vwapi" {
		key = "VWID"
	}
	params := map[string][]string{
		key:       []string{this.merchant.Key},
		"OrderId": []string{order.OrderId},
		"Amount":  []string{order.Amount}}

	if this.apiType == "vwapi" || this.apiType == "apim" {
		params["Password"] = []string{this.merchant.Password}
	}
	err = sendRequestAndMap(&orderResp, params, this, "Unblock")
	return
}

func (this OrderManager) Refund(order Payment) (orderResp OrderResponse, err error) {
	var params map[string][]string
	if this.apiType == "vwapi" {
		params = map[string][]string{
			"VWID": []string{this.merchant.Key},
			"DATA": []string{fmt.Sprintf("OrderId=%s;Password=%s;Amount=%s", order.OrderId, this.merchant.Password, order.Amount)}}
	} else {
		params = map[string][]string{
			"Key":      []string{this.merchant.Key},
			"OrderId":  []string{order.OrderId},
			"Amount":   []string{order.Amount},
			"Password": []string{this.merchant.Password}}
	}
	err = sendRequestAndMap(&orderResp, params, this, "Refund")
	return
}

func (this OrderManager) Charge(order Payment) (orderResp OrderResponse, err error) {
	key := "Key"
	if this.apiType == "vwapi" {
		key = "VWID"
	}
	params := map[string][]string{
		key:       []string{this.merchant.Key},
		"OrderId": []string{order.OrderId}}
	if this.apiType == "vwapi" || this.apiType == "apim" {
		params["Password"] = []string{this.merchant.Password}
		params["Amount"] = []string{order.Amount}
	}
	err = sendRequestAndMap(&orderResp, params, this, "Charge")
	return
}

func (this OrderManager) PayStatus(order Payment) (orderResp OrderResponse, err error) {
	params := map[string][]string{
		"Key":     []string{this.merchant.Key},
		"OrderId": []string{order.OrderId}}

	if this.apiType == "vwapi" {
		params = map[string][]string{
			"VWID": []string{this.merchant.Key},
			"DATA": []string{fmt.Sprintf("OrderId=%s", order.OrderId)}}
	}
	err = sendRequestAndMap(&orderResp, params, this, "PayStatus")
	return
}

func (this OrderManager) GetState(order Payment) (orderResp OrderResponse, err error) {
	params := map[string][]string{
		"Key":     []string{this.merchant.Key},
		"OrderId": []string{order.OrderId}}
	err = sendRequestAndMap(&orderResp, params, this, "GetState")
	return
}
