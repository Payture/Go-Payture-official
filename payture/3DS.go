package payture

import (
	"net/http"
)

func (this APIManager) APIThreeDSPAY(orderId string, paRes string) (*http.Response, error) {
	return this.send3DS("Pay3DS", orderId, paRes)
}

func (this APIManager) APIThreeDSBlock(orderId string, paRes string) (*http.Response, error) {
	return this.send3DS("Block3DS", orderId, paRes)
}

func (this APIManager) send3DS(method string, orderId string, paRes string) (*http.Response, error) {
	params := map[string][]string{
		"Key":     []string{this.merchant.Key},
		"OrderId": []string{orderId},
		"PaRes":   []string{paRes}}
	return sendRequest(this, method, params)
}

func (this EwalletManager) PaySubmit3DS(md string, paRes string) (*http.Response, error) {
	params := map[string][]string{
		"MD":    []string{md},
		"PaRes": []string{paRes}}
	return sendRequest(this, "PaySubmit3DS", params)
}
